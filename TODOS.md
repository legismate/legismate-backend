# things we should probably do 
 - have address be on the context of all requests, with helper methods on an `Address` struct that will return level values
   ```golang
   // this will return "seattle" for levelEnum=City, etc 
   (A *Address) func getLevelValue(level LevelEnum) {}
   ```  