# Application Instructions

OS: Linux

All instructions bellow asume you are in directory translatorapi

## Starting the Database
To start the database, run the following command:
```sh
sudo docker-compose up
```

## Starting the Server
Run the server using:
```sh
go run server.go
```

## Accessing the Application
Once the server is running, open the following URL in your browser:
```
http://localhost:8080/
```
Now you can use the application.

---



## Testing

For tresting you need to initialize new database, run docker by:

```sh
sudo docker run --name postgres-test -e POSTGRES_USER=test_user -e POSTGRES_PASSWORD=test_pass -e POSTGRES_DB=translatorapi_test -p 5433:5432 -d postgres

```

After that you need to test by typing:
```sh
go test ./resolvers_test.go
```

To end image of test sql:
```sh
sudo docker rm postgres-test
```
## Example Queries and Mutations

### Creating a Word
#### Request:
```graphql
mutation {
  createWord(polishWord: "a") {
    polishWord
    id
  }
}
```
#### Response:
```json
{
  "data": {
    "createWord": {
      "polishWord": "a",
      "id": "1"
    }
  }
}
```

### Creating a Word with Translation and Sentence
#### Request:
```graphql
mutation {
  createWord(polishWord: "Z", englishWord: "b", sentence:"c") {
    id
    polishWord
  }
}
```
#### Response:
```json
{
  "data": {
    "createWord": {
      "id": "2",
      "polishWord": "Z"
    }
  }
}
```

### Creating a Translation
#### Request:
```graphql
mutation {
  createTranslation(polishWord: "a", englishWord:"b") {
    id
    wordID
    englishWord
    examples {
      id
      translationID
      sentence
    }
  }
}
```
#### Response:
```json
{
  "data": {
    "createTranslation": {
      "id": "2",
      "wordID": "1",
      "englishWord": "b",
      "examples": []
    }
  }
}
```

### Creating an Example Sentence
#### Request:
```graphql
mutation {
  createExample(polishWord: "a", englishWord: "b", sentence: "bbbb") {
    id
    translationID
    sentence
  }
}
```
#### Response:
```json
{
  "data": {
    "createExample": {
      "id": "2",
      "translationID": "2",
      "sentence": "bbbb"
    }
  }
}
```

---

## Queries

### Retrieving Words and Their Translations
#### Request:
```graphql
query {
  words {
    polishWord
    id
    translations {
      englishWord
      examples {
        id
        sentence
      }
    }
  }
}
```
#### Response:
```json
{
  "data": {
    "words": [
      {
        "polishWord": "a",
        "id": "1",
        "translations": [
          {
            "englishWord": "b",
            "examples": [
              {
                "id": "2",
                "sentence": "bbbb"
              }
            ]
          }
        ]
      },
      {
        "polishWord": "Z",
        "id": "2",
        "translations": [
          {
            "englishWord": "b",
            "examples": [
              {
                "id": "1",
                "sentence": "c"
              }
            ]
          }
        ]
      }
    ]
  }
}
```

### Retrieving Translations for a Word
#### Request:
```graphql
query {
  translations(polishWord:"a") {
    id
    wordID
    englishWord
    examples {
      sentence
    }
  }
}
```
#### Response:
```json
{
  "data": {
    "translations": [
      {
        "id": "2",
        "wordID": "1",
        "englishWord": "b",
        "examples": [
          {
            "sentence": "bbbb"
          }
        ]
      }
    ]
  }
}
```

---

## Deleting Data

### Deleting a Word
#### Request:
```graphql
mutation {
  deleteWord(polishWord:"Z")
}
```
#### Response:
```json
{
  "data": {
    "deleteWord": true
  }
}
```

### Deleting a Translation
#### Request:
```graphql
mutation {
  deleteTranslation(polishWord: "a", englishWord: "b")
}
```
#### Response:
```json
{
  "data": {
    "deleteTranslation": true
  }
}
```

---

## Verifying Data After Deletion

### Retrieving Words
#### Request:
```graphql
query {
  words {
    polishWord
    id
    translations {
      englishWord
      examples {
        id
        sentence
      }
    }
  }
}
```
#### Response:
```json
{
  "data": {
    "words": [
      {
        "polishWord": "a",
        "id": "1",
        "translations": []
      }
    ]
  }
}
```


## Replacing The Translation

### Adding trnaslation

#### Request:

```graphql
mutation{
  createTranslation(polishWord: "a", englishWord: "b"){
    wordID
    englishWord
   }
 }

```

#### Response:

```json
{
  "data": {
    "createTranslation": {
      "wordID": "1",
      "englishWord": "b"
    }
  }
}
```

### Replacing translation:


### Request:
```graphql
mutation{
  replaceTranslation(polish_word: "a", englishWord: "b", newTranslation: "c"){
    id
    englishWord
    
  }
}
```


### Response:
```json
{
  "data": {
    "replaceTranslation": {
      "id": "4",
      "englishWord": "c"
    }
  }
}
```



### Check:

### Request: 
```graphql
query {
  words {
    polishWord
    id
    translations {
      englishWord
      examples {
        id
        sentence
      }
    }
  }
}
```
### Response:
```json
{
  "data": {
    "words": [
      {
        "polishWord": "a",
        "id": "1",
        "translations": [
          {
            "englishWord": "c",
            "examples": []
          }
        ]
      }
    ]
  }
}
```