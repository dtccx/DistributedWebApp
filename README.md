# DistributedWebApp
# Group members
## Changxing Cao
## Xinglun Xu

# How to Install the web application
1.  Install APIs we used:           
    go get github.com/gorilla/sessions    

2.  Run the server:           
    go run server.go

3.  enter the website:        
    http://localhost:8080/index.html is the enter of the website.

# Idea       
Time to create a fantastic product that will make you gazillions!!    
I think no one has ever come up with this idea before, so don't share it as we will be barraged with legal cases later on if you do.    
People are dying to share every little mindless thing they do throughout their lives, throughout their days, every minute. But given all that sharing, there isn't much time for actually saying much that's meaningful, so we can limit what they will say to a nice small fixed size, say 100 characters.   
The first stage of the overall project will be a simple web application, comprised of html and Go on a single machine. To further simplify matters, for this stage, you do not have to keep the information in files, instead just keepieverything in in-memory data structures. Note we are not using a database for this application.   
The second stage will involve splitting off the back end off to a separate server machine, implemented in Go and communication done via RPC and making sure to use goroutines.    
And finally, in the third stage, to improve reliability, we will replicate the server.    
What are the requirements for this application? You are designing the system, so aside from the description   
written above, much of the way the system works will be up to you. Clearly, there must be:    
Users who log in. They will need passwords. They may want to cancel their accounts. The ability to write / send the messages that the system is built for.    

## The major components will be graded as
1. Basic web app: 30%, due date: April 8th
2. Separating off the back end and making sure it’s multithreaded: 30%, due date: April 22nd
3. Replicated back end: 30%, due date: May 6th

# Some dramatic bugs in Golang
1.  In datasture, the var's first word must be Uppercase; Otherwise Json or other function will not get them. : )
2.  like.js cannot be recognized by js file. :) rename it!!!

## Picture explaination for The Web each function:
1.  Login and sign up Main Page:              
![](img/img01.png)
2.  After Login, the home page will be :        
![](img/img02.png)  
3.  Cancel their account function:              
    (delete account after their login)                  
![](img/img03.jpeg)       
![](img/img04.png)    
4.   Write and send msg:            
![](img/img05.png)         
![](img/img06.png)     
5.  specify whose messages:                 
    For each message, we show the name of the user who sent the message;              
![](img/img07.png)    
6.  Like Function:
    Click on the message, you will like this messages

7.  Show all the message you Like       

8.  Show the message's like num (how many people liked this message)
