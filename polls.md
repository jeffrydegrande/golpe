FORMAT: 1A

# Polls
Polls is a simple API allowing consumers to view polls and vote in them.

# Group Questions
Resources related to questions in the API.


## Question [/questions/{id}]
A Question contains information about a Poll.

+ Parameters

    + id: 1 (number, optional) - The ID of the desired question
        + Default: `100`

+ Attributes (Question Base)
    + id: 250FF (string)
    + created: 1415203908 (number) - Time stamp

### Retrieve a Question [GET]
Retrieves the Question with the given ID

+ Response 200 (application/json)
        [
            {
                "question": "Favourite programming language?",
                "published_at": "2014-11-11T08:40:51.620Z",
                "url": "/questions/1",
                "choices": [
                    {
                        "choice": "Swift",
                        "url": "/questions/1/choices/1",
                        "votes": 2048
                    }, {
                        "choice": "Python",
                        "url": "/questions/1/choices/2",
                        "votes": 1024
                    }, {
                        "choice": "Objective-C",
                        "url": "/questions/1/choices/3",
                        "votes": 512
                    }, {
                        "choice": "Ruby",
                        "url": "/questions/1/choices/4",
                        "votes": 256
                    }
                ]
            }
        ]


## Questions [/questions]

+ Attributes (array[Question])

### List All Questions [GET]
Returns a list of your Questions.

+ Response 200 (application/json)

        [
            {
                "question": "Favourite programming language?",
                "published_at": "2014-11-11T08:40:51.620Z",
                "url": "/questions/1",
                "choices": [
                    {
                        "choice": "Swift",
                        "url": "/questions/1/choices/1",
                        "votes": 2048
                    }, {
                        "choice": "Python",
                        "url": "/questions/1/choices/2",
                        "votes": 1024
                    }, {
                        "choice": "Objective-C",
                        "url": "/questions/1/choices/3",
                        "votes": 512
                    }, {
                        "choice": "Ruby",
                        "url": "/questions/1/choices/4",
                        "votes": 256
                    }
                ]
            }
        ]

### Create a New Question [POST]

You may create your own question using this action. It takes a JSON object
containing a question and a collection of answers in the form of choices.

+ Attributes (Question Base)
+ Request (application/json)
+ Response 201 (application/json)

    + Attributes (Question)

    + Headers

            Location: /questions/2

# Data Structures

## Question Base (object)

+ question (string) - The question
+ choices (array[string]) - A collection of choices.