How to run your server?
main.go is located in the next-server folder due to some issues with codespace
-The server can be run by entering go run main.go --port=4000

-Then you should see Server listening on :4000

Which functionality was the most educationally enriching?

-The inactivity timeout and selecting-basing Handling.

- Implementing the 30-second inactivity timeout using Go’s select statement with channels and timers was the most conceptually enriching. It forced you to think asynchronously, managing:

Real-time input from the client (msgChan)

Timeout triggers (time.Timer)

Clean exits (quitChan)

Which functionally required that you do the most research?

-The functionally that required the most research was the safe logging per client with sync.Mutex.

Demo video link
https://youtu.be/XffY1RRhgi0

short video explanation link
https://youtu.be/fT_MUVR8BDo