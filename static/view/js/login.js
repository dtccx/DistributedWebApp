var url = "http://localhost:8080";

$( "#btnReg" ).click(function() {
  alert("btnReg")
  $.get( url+"/User/Register",{userName:$("#userName").val() ,password:$("#password").val()}).done(function( data ) {
      console.log(data);
    	if(data=="0") alert("User Name is occupied, try another one.");
    	else{
    		console.log("Start replace");
    		window.location.replace(url);
    		console.log("End replace");
    	}
  	});
});

$("#btnLogin").click(function(){
	$.get( url+"/User/Login",{userName:$("#userName").val() ,password:$("#password").val()}).done(function( data ) {
      console.log("yaayay");
      var loginSuccess = data==true;
  		if(loginSuccess){
  			window.location.replace(url);
  		}
  		else{
  			alert("Either the Username or the password is incorrect");
  		}
  	});
});
