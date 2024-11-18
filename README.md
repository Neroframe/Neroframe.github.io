 # Executive Summary

In this project, I developed a healthcare CRUD (Create, Read, Update, Delete) application using Go (Golang) and PostgreSQL, deployed on the Railway platform. The primary objective was to build a functional web application that allows users to perform CRUD operations on healthcare-related data. The backend was most annoying part, because each table required basic CRUD operations and respective web pages. Those, I spent lots of hours just copy/pasting and adjusting base code to meet each tables requirements.

### Completed Tasks:

Database Design and Implementation: Established a PostgreSQL database with relevant tables (11 in total) such as users, patients, doctors, disease, and publicservant. I ensured proper relationships and constraints between tables to maintain data integrity.

Backend Development: Developed the backend using Go, implementing handlers for various endpoints to perform CRUD operations. This includes inserting new records, fetching existing data, updating records, and deleting entries from the database.

Database Connection Management: Configured the application to connect securely to the PostgreSQL database hosted on Railway. Addressed challenges related to environment variables and ensured that the DATABASE_URL was correctly utilized in different environments (development and production).

Deployment: Successfully deployed the application on Railway, setting up the necessary environment configurations. Linked the web service with the PostgreSQL database within the Railway project for seamless integration.

Testing and Debugging: Conducted thorough testing to ensure all routes and functionalities work as expected. Resolved issues related to database connectivity, foreign key constraints, and SSL requirements.

### Incomplete Tasks:

Authorization and Security Measures: I did not implement user authorization, authentication, or other security features. As per the project's scope, these aspects were not objectives for this phase of development. I believe that there are LOTS OF ways to break my system. While I tried to address possible inputs, I cannot be fully certain that the system would not fail under extreme edge cases.

Frontend Development: The project focused primarily on the backend and database interactions. Any user interface components were minimal or not fully developed.
