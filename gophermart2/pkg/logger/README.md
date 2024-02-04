## logger - logger for app
### Init
Logger initialization

### Fatal
Starts a new message with fatal level. The os.Exit(1) function is called by the Msg method, which terminates the program immediately.


### Error
Error starts a new message with error level. 
You must call Msg on the returned event in order to send the event.


### Info
Info starts a new message with info level.
You must call Msg on the returned event in order to send the event.

### Debug
Debug starts a new message with debug level.
You must call Msg on the returned event in order to send the event.