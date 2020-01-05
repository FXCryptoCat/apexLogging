# apexLogging
Go program that pulls data from a Neptune Apex and pushes to InfluxDb.

# Note
Please let me know if you have any feature requests. At this point it is doing what I need
I would love to have some ideas.

If you know how to login to apex fusion from go, please let me know. I would like to 
add tracking for test results.

# Apex Interface
The protocol to the APEX Classic Console was reverse engineered and used in this project.

## Security
Update: Now you only need to provide the username and password to the apexLogger, it will login
and get the cookie for you. 

The security for the APEX is minimal. All you need to do is get the vale of a cookie and send 
it with your HTTP requests.

### Apex Interfaces Used
All of this was reverse engineered. 

>  http://111.123.12.15/rest/status/?_=[MILLIS]
>> Millis are the time at which you are requesting status. Usually current time.  
>> This call is requested when loading the main page. We care about the Inputs and Outputs.
>> See the ApexStatus struct in apex_status.go for more details   

  
>  "http://111.123.12.15/rest/[LOGTYPE]?days=[DAY]&sdate=[DATE]&_=[CURRENT_TIME_SECONDS]"
>> LOGTYPE is ilog(input) or olog(output). Others supported by apex are dlog (dose) and tlog(??)  
>> DAY is the number of days of logs to retrieve.  
>> DATE  YYMMDD for the first day to retrieve logs  (FROM TIME)
>> CURRENT_TIME_SECONDS, the time to stop getting logs. (TO TIME)
>> This called is used when loading the 'Graph' page. It brings in everything.
>> See apex_log.go, apex_input_log.go and apex_output.go for details


# Installation
The apexLogger can install itself as a service on a RaspberryPi 3b+. Let me know if you 
try this on another OS and it works so I can update the README. 
