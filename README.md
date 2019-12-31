# apexLogging
Go program that pulls data from a Neptune Apex and pushes to InfluxDb.

#Apex Interface
The protocol to the APEX Classic Console was reverse engineered and used in this project.

## Security
The security for the APEX is minimal. All you need to do is get the vale of a cookie and send 
it with your HTTP requests.

###How to Get the Cookie
Call Cookie Monster. If he isn't available use the following steps.

- Open Chrome and navigate to your APEX classic dashboard
- Sign in
- Open Chrome's settings
- Navigate to the "Advanced Settings"
- Go to the "Site Settings"
- Go to "Cookies and Site Data"
- Go to "See all cookies and site data"
- Search for your APEX's IP address
- You will find a couple of cookies, look for **connect.sid**. Use this value as your cookie.
> I'm not sure what happens when/if it expires. If some can get me the code to re-authorize I'll be happy  
> to add you as a contributor and review it. 


###Apex Interfaces Used
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


