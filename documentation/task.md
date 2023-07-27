# Understanding the problem
Imagine you’re working at an AI startup as an analyst. You are tasked with creating retention charts so we can get our investors up to speed on how people are using our product.

The usage history is a recording of inferences made on the website.

Draw a retention chart, which shows the month-on-month retention of users, based on whether they ran an inference on a particular month or not.

### Querying the data
It's a 50MB JSON file, with some 202,757 records. I feel like before I can start plotting the information given here, I should get an understanding of the shape of the data I'm dealing with.

In order to query it, I decided to use [OctoSQL](https://github.com/cube2222/octosql). OctoSQL allows you to query varies data formats, including JSON and CSV.

It supports JSON in the form of [JSONLines](https://jsonlines.org), which felt like a convenient form to read the site logs provided.

```bash
$ jq --compact-output '.[]' \
	gooey_inferences.json > inferencelines.json
```

From `man jq`

```plain
--compact-output / -c:

By default, jq pretty-prints JSON output. 

Using this option will result in more compact output by 
instead putting each JSON object on a single line.
```


### Initial impressions

#### Schema

```JSON
{
  "is_anonymous": true, 
  // boolean
  "recipe": "DeforumSD", 
  // string
  "run_id": "dc158e8add22da7690cd35efe6f1c0ef6a43932868b8958c316a85c95f34e612", 
  // string, length 64, probably hash
  "user_id": "e02d31832fdba9afd55b6e2107e097124b7ddc28d6e6eac524f5afb42d19cc6b", 
  // string, length 64, probably hash
  "timestamp": "2023-03-05T21:17:17.679325" 
  // ISO8601 timestamp with local time — no timezone provided
}
```


#### The timestamps are in random order

```plain
"select timestamp from inferencelines.json limit 10"
+------------------------------+
|          timestamp           |
+------------------------------+
| '2023-03-05T21:17:17.679325' |
| '2023-03-05T21:17:57.277849' |
| '2023-03-05T21:18:00.957005' |
| '2023-03-05T21:18:05.409089' |
| '2023-03-05T21:18:24.693359' |
| '2023-05-03T12:54:13.401587' |
| '2023-05-03T12:54:20.963993' |
| '2023-06-05T19:14:54.371972' |
| '2023-06-05T19:40:11.781434' |
| '2023-06-20T17:33:55.982841' |
+------------------------------+
```

#### The timestamps have a huge range

```plain
# DESCENDING

"select timestamp from inferencelines.json order by timestamp desc limit 5"
+------------------------------+
|          timestamp           |
+------------------------------+
| '2023-06-20T23:57:36.313687' |
| '2023-06-20T23:55:00.079252' |
| '2023-06-20T23:53:35.331227' |
| '2023-06-20T23:53:11.789879' |
| '2023-06-20T23:52:34.112602' |

# ASCENDING

"select timestamp from inferencelines.json order by timestamp asc limit 5"
+-----------------------+
|       timestamp       |
+-----------------------+
| '1970-01-01T00:00:00' |
| '1970-01-01T00:00:00' |
| '1970-01-01T00:00:00' |
| '1970-01-01T00:00:00' |
| '1970-01-01T00:00:00' |
+-----------------------+
```

#### There is some noisy data — some timestamps are at unix epoch.
- 2783 out of 202,757 timestamps can be ignored


```plain
"select count(*) from inferencelines.json where timestamp='1970-01-01T00:00:00'"

+-------+
| count |
+-------+
|  2783 |
+-------+

"select count(*) from inferencelines.json"
+--------+
| count  |
+--------+
| 202757 |
+--------+
```

#### The timestamps after unix epoch are reliable
```plain
$ octosql "select timestamp from inferencelines.json order by timestamp asc limit 2800" | tail -30

| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '1970-01-01T00:00:00'        |
| '2022-12-21T15:13:03.179443' |
| '2022-12-21T15:22:30.824830' |
| '2022-12-21T15:23:13.237637' |
| '2022-12-21T15:32:11.241788' |
| '2022-12-21T15:34:24.954346' |
| '2022-12-21T15:35:38.757209' |
| '2022-12-21T15:37:18.792011' |
| '2022-12-21T15:37:56.250975' |
| '2022-12-21T15:47:25.154315' |
| '2022-12-21T15:48:41.357170' |
| '2022-12-21T15:51:57.291146' |
| '2022-12-21T15:52:23.585429' |
| '2022-12-21T15:53:31.977901' |
| '2022-12-21T16:05:53.842276' |
| '2022-12-21T16:08:53.690751' |
| '2022-12-21T16:13:25.072171' |
| '2022-12-21T16:14:18.323783' |
+------------------------------+
```

- Legitimate range of months is from Dec 2022 to Jun 2023: 7 months

#### Recipes are in such a number that retention *can* be analysed  individually
```plain
"select count(distinct recipe) FROM inferencelines.json"
+-----------------------+
| count_distinct_recipe |
+-----------------------+
|                    28 |
+-----------------------+
```


# Ingesting the data
## Choosing a data store
In terms of finding distinct userIDs for each month, to form the basis of retention chart, I’d like to store in something more query-able that plain JSON.

Seeing how this a batch processing job, and the schema is stable across the file, I think ingesting the data into a relational database. 

With the intention of making the process of running the code easy, I’ll store the data in SQLite, and query that to generate datasets for chart-making.

Storing in relational data store makes querying easier, and makes future analyses easier.

Note: “run_id” is not unique, as I encountered when entering the data into SQLite.

```JSON
{
	"is_anonymous":false,
	"recipe":"EmailFaceInpainting#2",
	"run_id":"7fca2ecd62a553a144bf2ccc6684886a32a534de71741ef936e3dbfcd6a9208e",
	"user_id":"62ba6228f1d7918283cd7517ab72bd7e5ae02560b6ab7db74b531283f0b776e0",
	"timestamp":"2022-12-22T05:14:24.610073"
}

{
	"is_anonymous":false,
	"recipe":"SEOSummary",
	"run_id":"7fca2ecd62a553a144bf2ccc6684886a32a534de71741ef936e3dbfcd6a9208e",
	"user_id":"62ba6228f1d7918283cd7517ab72bd7e5ae02560b6ab7db74b531283f0b776e0",
	"timestamp":"2022-12-26T10:55:25.856325"
}
```