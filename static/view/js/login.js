var url = "http://localhost:8080";

$( "#btnReg" ).click(function() {
  alert("btnReg")
  $.get( url+"/User/Register",{user:$("#userName").val() ,password:$("#password").val()}).done(function( data ) {
      console.log(data);
    	if(data=="0") alert("User Name is occupied, try another one.");
    	else{
    		console.log("Start replace");
        alert("Sign Up Successfully, you can login now");
    		window.location.replace(url + "/index.html");
    		console.log("End replace");
    	}
  	});
});

$("#btnLogin").click(function(){
	$.get( url+"/User/Login",{user:$("#userName").val() ,password:$("#password").val()}).done(function( data ) {
      console.log("yaayay");
      console.log(data);
      var loginSuccess = data=="true";
  		if(loginSuccess){
        //to the home page
  			window.location.replace(url + "/index.html");
  		}
  		else{
  			alert("Either the Username or the password is incorrect");
  		}
  	});
});
