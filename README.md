# file-store
Simple file store service 

This is a simple Go HTTP server that provides endpoints for storing, listing, removing, updating, counting words, and finding frequent words in uploaded files.

# Setup:-
1. Clone this repository:
- git clone https://github.com/nehaLohia27/file-store.git

2. Navigate to repository
- cd file-store

3. Run the server.go
- go run server.go

# Endpoints:

1. Add File
**POST /store/add

- Add one or more files to the server's storage. Each file should be sent as a multipart form data with the key "file". Files are stored in the uploads/ directory.

2. List Files
GET /store/ls

- List all files stored on the server.

3. Remove File
DELETE /store/rm?file=<file_name>

- Remove a file from the server's storage. Provide the filename as a query parameter.

4. Update File
PUT /store/update?file=<file_name>

- Update a file in the server's storage. Provide the filename as a query parameter and send the new file as multipart form data with the key "file".

5. Word Count
GET /store/wc

- Count the total number of words across all uploaded files.

6. Frequent Words
GET /store/freq-words?limit=<limit>&order=<asc|desc>

- List the most frequent words in uploaded files. You can specify the limit and order of the results.


# Usage Example: 

1. Add a file:
./store add file1.txt

2. Update the file 
./store update file1.txt

3. list the files
./store ls

4. Remove a file 
./store rm file1.txt

5. Count the total words in a file
./store wc

6. Print the most frequent words
./store freq-words






