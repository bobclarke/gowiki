const { events, Job, Group } = require("brigadier");


events.on("exec", function(e, project) {
  console.log("Hello world");
})

events.on("push", function(e, project){
    console.log("Push event recieved from github")
})