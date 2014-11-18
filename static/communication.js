$( document ).ready(function() {
    var c=new WebSocket('ws://localhost:3000/sock');
    c.onopen = function(){
      send({Id: "0", Message: "", Type: "connect"})
      c.onmessage = function(response){
        var data = JSON.parse(response.data);
        console.log(data);
        switch (data.Type){
            case "get_value":
                var selector = data.Message
                var id = data.Id
                val = $(selector).val()
                send({Id: id, Message: val, Type: "reply"})
            break
            case "ping":
                send({Id: id, Message: "pong", Type: "reply"})
            break
            case "get_content":
                var selector = data.Message
                var id = data.Id
                val = $(selector).html()
                send({Id: id, Message: val, Type: "reply"})
            break
            case "set_content":
                var selector = data.Message.Selector
                var id = data.Id
                val = $(selector).html(data.Message.Content)
            break
            case "set_value":
                var selector = data.Message.Selector
                var id = data.Id
                val = $(selector).val(data.Message.Content)
            break
            
        }
      };
    }
    
    function send(data){
        console.log("send", data)
        c.send(JSON.stringify(data))
    }
    
})