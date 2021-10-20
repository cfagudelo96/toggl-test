# toggl-test

To run the project, execute the following commands:

`go mod download`
`go run .`

The application will be occupy by default in the port 3000, but it can be changed
by changing the environment variable PORT to a different one.

To execute the unit test the command `go test ./...`

Only unit test for the deck domain package were provided for time constraints.

## Create new deck

To create a new deck the following endpoint must be consumed:

`POST <host>/v1/decks`

To specify the cards in the deck, the query parameter `cards` must be specified with the list
of card codes to add separated by commas.

If the deck to create should not be shuffled, the query parameter `shuffled` must be provided with the value `n`. If
the query parameter is not specified or has a different value, the deck will be shuffled.

The following command shows how to use the endpoint running the app locally, creating an unshuffled deck of 5 cards:

`curl --location --request POST 'http://localhost:3000/v1/decks?cards=AS,KD,AC,2C,KH&shuffled=n'`

## Open a deck

To open a deck the following endpoint must be consumed:

`GET <host>/v1/decks/<Deck ID>`

The following command shows how to use the endpoint running the app locally with an example deck:

`curl --location --request GET 'http://localhost:3000/v1/decks/7bd85b70-5662-4022-ad33-e916359df19e'`

## Draw cards

To draw cards from a deck the following endpoint must be consumed:

`POST <host>/v1/decks/<Deck ID>/draw`

To specify the amount of cards to draw, a JSON body should be provided as follows:

```json
{
  "amount": 3
}
```

If the amount of cards to be drawn is greater than the amount of cards in the deck, all the cards
will be drawn and no error will happen (this was a design decision).

The following command shows how to use the endpoint running the app locally with an example deck:

`curl --location --request POST 'http://localhost:3000/v1/decks/43cc860b-f74f-4421-8858-6f14c2f1c476/draw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 2
}'`
