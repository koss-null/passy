# Passy

Passy is a tool to create and store your encripted passwords in your git repo.  

## Fast start

Install the tool:  
`go install github.com/koss-null/passy@latest`

See help page:  
`passy -h`

Generate some passwords:  
`passy -c --readable` - generates some readable password
`passy -c --safe` - generates some safe, kind of readable password
`passy -c --insane` - generates some really safe password

## TBD:

- **Password Storage** Add ability to store passwords in key-value format to be able to get the pass by key
- Configuring and storing your encoded passwords in git repo 
- Interactial mode

