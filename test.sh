#!/bin/bash

# Configuration
URL="http://localhost:8080/api/v1/is-allowed"
TOTAL_REQUESTS=40
DELAY=0.1 # Seconds between requests

echo "Starting rate limit test against $URL"
echo "Sending $TOTAL_REQUESTS requests with a ${DELAY}s delay..."
echo "--------------------------------------------------------"

for ((i=1; i<=TOTAL_REQUESTS; i++)); do
    # Execute curl and capture the HTTP status code and body
    # -s: Silent mode
    # -w: Write out specific format (status_code)
    RESPONSE=$(curl -s -w "\n%{http_code}" "$URL")
    
    # Split the response into body and status
    BODY=$(echo "$RESPONSE" | head -n 1)
    STATUS=$(echo "$RESPONSE" | tail -n 1)

    if [ "$STATUS" -eq 200 ]; then
        echo "Request #$i: [OK 200] - $BODY"
    elif [ "$STATUS" -eq 429 ]; then
        echo "Request #$i: [LMT 429] - $BODY"
    else
        echo "Request #$i: [ERR $STATUS] - $BODY"
    fi

    sleep $DELAY
done

echo "--------------------------------------------------------"
echo "Test Complete."
