# DistributedWebApp

Time to create a fantastic product that will make you gazillions!!    
I think no one has ever come up with this idea before, so don't share it as we will be barraged with legal cases later on if you do.    
People are dying to share every little mindless thing they do throughout their lives, throughout their days, every minute. But given all that sharing, there isn't much time for actually saying much that's meaningful, so we can limit what they will say to a nice small fixed size, say 100 characters.   
The first stage of the overall project will be a simple web application, comprised of html and Go on a single machine. To further simplify matters, for this stage, you do not have to keep the information in files, instead just keepieverything in in-memory data structures. Note we are not using a database for this application.   
The second stage will involve splitting off the back end off to a separate server machine, implemented in Go and communication done via RPC and making sure to use goroutines.    
And finally, in the third stage, to improve reliability, we will replicate the server.    
What are the requirements for this application? You are designing the system, so aside from the description   
written above, much of the way the system works will be up to you. Clearly, there must be:    
Users who log in. They will need passwords. They may want to cancel their accounts. The ability to write / send the messages that the system is built for.    
The ability to specify whose messages you are interested in.    
Other features? If I have omitted something critical, let me know. If there are more features you want to add, feel free to expand.   
## The major components will be graded as
1. Basic web app: 30%, due date: April 8th
2. Separating off the back end and making sure it’s multithreaded: 30%, due date: April 22nd
3. Replicated back end: 30%, due date: May 6th
4. Product Demo: 10%, due date: May 10 in class.

## Deliverables:
1. Code and unit Tests in Github. We will check the Github for commits from all the team members and deduct points for people that don’t have any commits. Also if you copy from anyone you will get a 0 for the project.
2. 10-15 minute Demo to the instructor and TAs during Class.
