# API Gateway and Load Balancer

* Load Balancer - Distributes incoming traffic across multiple API Gateway instances, can be scaled horizontally by adding more instances to handle increased traffic. Elastic load balancers (e.g., AWS ELB) can automatically scale based on traffic load and ensure high availability and fault tolerance by distributing requests across multiple servers.

* API Gateway - Manages incoming HTTP requests, handling authentication, rate limiting, and routing. 

# Web Server

* Handle HTTP requests for image uploads, generate upload links, and fetch images. The stateless design allows horizontal scaling by adding more web server instances and can use container orchestration tools like Kubernetes to ensure high availability.

# Message Queue

* Acts as a buffer to decouple web servers from processing layers, ensuring reliable and scalable message delivery.

# Real-Time Processing

* Processes messages in real-time for tasks such as data cleaning, aggregation, and generating immediate predictions for uploaded images. This systems are usually designed to handle high-throughput data streams with low latency.

# Persistence

* Stores results from real-time and batch processing for quick access by clients.

# Batch Processing

* Processes data in batches for training and updating predictive models, running periodic reports, and complex data transformations. Can scale horizontally by adding more processing nodes.

# Model Serving

* Provides APIs to serve trained models and generate predictions, integrating with both real-time and batch processing layers. Can scale horizontally by deploying multiple instances of model servers.

# Questions for better decision making

Data Volume and Velocity:

What is the expected volume of image uploads per second?
What is the peak traffic we need to handle?

Processing Requirements:

What kind of data processing needs to be done in real-time versus batch?
What are the latency requirements for real-time processing?

Storage Requirements:

How much data do we need to store?
What are the retention policies for raw and processed data?

Scalability and Fault Tolerance:

How can we ensure that the system scales horizontally?
What are the backup and recovery strategies?

Security and Compliance:

What are the security requirements for data at rest and in transit?
Are there any compliance requirements (e.g., GDPR)?

Cost Management:

What is the budget for infrastructure and operations?
How can we optimize costs while ensuring performance?