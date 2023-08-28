# Exp Blog - Golang

### Project Description:
The Golang backend system for the multi-vendor blog application aims to create a robust and scalable platform that facilitates multiple vendors to publish and manage their blogs efficiently. The system will provide a user-friendly interface for vendors to register, create, edit, and manage their blog posts, catering to various topics and niches.


## Key Features:
- ### User Authentication and Authorization: 
    The system will have a secure user authentication mechanism, allowing vendors to register and log in to their accounts. Role-based authorization will be implemented to control access to specific functionalities based on user roles.

- ### Vendor Dashboard:
  Vendors will have access to a personalized dashboard where they can manage their profile, view their published blogs, draft new articles, track blog statistics, and handle comments and interactions.

- ### Blog Creation and Editing:
  Vendors can create new blog posts using a user-friendly editor, allowing them to add text, images, multimedia content, and tags to categorize their posts effectively. They can save drafts, schedule posts for future publication, or publish them immediately.

- ### Tagging and Categorization:
  The system will support tagging and categorization of blog posts, enabling users to discover content easily based on their interests.

- ### User Interaction:
  Readers can interact with blog posts by liking, disliking, and commenting on them. Vendors can moderate comments to maintain a healthy and engaging community.

- ### Search and Filtering:
  The application will offer robust search functionality, enabling users to search for blogs by keywords, tags, authors, or categories.

- ### Analytics and Insights:
  Vendors will have access to analytics and insights, providing them with valuable data on blog performance, reader engagement, and audience demographics.

- ### Notification System:
  The platform will implement a notification system to alert vendors about comments, likes, or other relevant activities related to their blogs.

- ### Vendor Collaboration:
  The system can support multiple authors collaborating on a single blog post, facilitating co-authored articles.

- ### Admin Panel:
  An admin panel will be provided to manage users, blog posts, comments, and overall system settings.

- ### APIs for Mobile and Frontend Integration:
  The backend will expose APIs, allowing easy integration with mobile apps and frontend applications.

Overall, the Golang backend system for the multi-vendor blog application will empower vendors to showcase their expertise through blogs, engage with their audience, and build a thriving community around diverse content. With an intuitive user interface, robust security measures, and powerful features, the platform aims to be a go-to choice for both vendors and readers seeking quality blog content.


## API Routes

### Authentication and User Management

- **Login**: `POST /api/v1/login`
- **Forgot Password**: `POST /api/v1/password/forgot`
- **Reset Password**: `POST /api/v1/password/reset`

### Users

- **Create User**: `POST /api/v1/users`
- **Get Users**: `GET /api/v1/users`
- **Get User by ID**: `GET /api/v1/users/:id`
- **Update User by ID**: `PUT /api/v1/users/:id`
- **Update User Avatar**: `PUT /api/v1/avatar/users/:id`
- **Delete User by ID**: `DELETE /api/v1/users/:id`

### Profiles

- **Create User Profile**: `POST /api/v1/profiles`
- **Get User Profiles**: `GET /api/v1/profiles`
- **Get User Profile by ID**: `GET /api/v1/profiles/:id`
- **Update User Profile by ID**: `PUT /api/v1/profiles/:id`
- **Update User Profile Picture**: `PUT /api/v1/avatar/profiles/:id`
- **Delete User Profile by ID**: `DELETE /api/v1/profiles/:id`

### Posts

- **Create Post**: `POST /api/v1/posts`
- **Get Posts**: `GET /api/v1/posts`
- **Get Post by ID**: `GET /api/v1/posts/:id`
- **Update Post by ID**: `PUT /api/v1/posts/:id`
- **Delete Post by ID**: `DELETE /api/v1/posts/:id`
- **Get User Profile Posts**: `GET /api/v1/user_posts/:id`

### Likes

- **Get Likes for Post**: `GET /api/v1/likes/:id`
- **Like Post**: `POST /api/v1/likes/:id`
- **Unlike Post**: `DELETE /api/v1/likes/:id`

### Comments

- **Create Comment for Post**: `POST /api/v1/comments/:id`
- **Get Comments for Post**: `GET /api/v1/comments/:id`
- **Update Comment by ID**: `PUT /api/v1/comments/:id`
- **Delete Comment by ID**: `DELETE /api/v1/comments/:id`

### Comment Replies

- **Create Comment Reply for Comment**: `POST /api/v1/comment/replyes/:id`
- **Get Comment Replies for Comment**: `GET /api/v1/comments/replyes/:id`
- **Update Comment Reply by ID**: `PUT /api/v1/comments/replyes/:id`
- **Delete Comment Reply by ID**: `DELETE /api/v1/comments/replyes/:id`
