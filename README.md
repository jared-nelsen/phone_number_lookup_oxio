# OXIO Phone Number Lookup API

This is a Go-based REST API service that provides phone number lookup functionality. It validates and parses phone numbers according to E.164 format specifications

## Notes

## a. Running solution

- I chose Golang for this exercise

  ## Running the Server

  1. Install dependencies:
     go mod tidy

  2. Run the server:
     go run main.go

  3. Server will start on 8080

  ## Building

  go build
  ./phone_number_lookup

  ## Run Tests

  go test

## b. Tech Choices

- I chose Golang for a few reasons:
  - 1. A phone number lookup/validation API is probably going to be really high scale and Golang scales well without a lot of tuning.
    - Its straightforward to get a really high Request Rate with Golang in this type of problem
  - 2. Fast execution + serialization
    - Although users of this API are predominantly affected by network latency using this API it still matters how much time it takes
      to serialize/deserialize the json and also do the string operations. Golang is great at IO bound work and pretty quick on execution
      so it seemed like a solid choice.
  - 3. Gin Gonic makes it really easy to stand up an API

## c. Proposed deployment

- Amazon ECS with simple horizontal scaling on memory. Golang is really efficient and we could easily hit really high load for this API with this model.

## d. Assumptions

- That message will only contain the prescribed characters in the string. There could be issues with validation unless I was more thorough.
- Basically more of the same for all different issues like content type and request validation

## e. Improvements

- Support full list of countries
- Consolidate Maps - Instead of two maps for the Country and Area code length I could have the Country code mapped to a struct with all relevant associated data.
  This would allow us to easily extend that struct with other metadata in the future without a proliferation of maps
- Add a flow diagram - I had to think a lot about how the validation flow should go. It would be good to include a representation of that for others to use.

## API Endpoint

### GET /v1/phone-numbers

Parses and validates a phone number, returning its components.

### Query Parameters

- `phoneNumber` (required): The phone number to parse

  - Format: E.164 ([+][country code][area code][local phone number])
  - The `+` is optional
  - Phone number must be a sequence of digits
  - Spaces are allowd between country code, area code, and local phone number
  - Any other characters are invalid
  - Any other space placement is invalid

- `countryCode` (optional): ISO 3166-1 alpha-2 country code
  - Required if te phone number doesn't include a country code
  - Must be exactly 2 characters
  - Examples: `US`, `MX`, `ES`

## Response Format

### Success Response (200 OK):

```json
{
  "phoneNumber": "+12125690123",
  "countryCode": "US",
  "areaCode": "212",
  "localPhoneNumber": "5690123"
}
```

### Error Response (400 Bad Request):

```json
{
  "phoneNumber": "631 311 8150",
  "error": {
    "countryCode": "required value is missing"
  }
}
```

## Architecture

Basically, there are two top level cases that define how we should try to parse and validate the number:

- The number has a country code in it
- The number does not have a country code in it

I basically split the execution paths on this fact and carry out the validation step by step in each case. I also split out a lot of the validation into individual functions.
The combination of these made it really easy to write unit tests.
