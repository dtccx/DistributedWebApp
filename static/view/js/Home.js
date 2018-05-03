var url = "http://localhost:8080";

var browseStack = []
var elements = []
var lastMsgId = 0
var currentOffset = 0

const NORMAL_MSG_TYPE = 0
const LASTEST_MSG_FLAG = -2

function playStackTop(){
	console.trace();
	console.log("playStackTop called");
	elements = [];
	if(browseStack.length==0){
		getFeeds();
	}
	else{
		var temp = browseStack[browseStack.length-1];
		temp.action(temp.params);
	}
	updateTrace();
}

function updateTrace(){
	var valS = "Home";
	for(i=0; i<browseStack.length; i++){
		valS = valS + ">>"+browseStack[i].params.pageTag;
	}
	$("#trace-content").html(valS);
}

function resetStack(){
	browseStack = [];
}


function closeWithRow(s){
	return "<div class=\"row\" id=\"home-table-row\">" + s + "</div>";
}

function createCard(id, row, cardType){
	var ret = "";
	var cardClass = "track-card";
	switch(cardType){
		case NORMAL_MSG_TYPE:
			cardClass = "normal-msg-card"
			ret = createMsgCard(id, row);
			break;
		default:
			break;
	}

	var temp = `<div class="col-sm-4">
				<div class="card home-table-card2 ${cardClass}">
					<div class="card-block home-table-card" id="card-${cardType}-${row.ID}">${ret}</div>
				</div>
			 </div>`;
	return temp;
}

function getMsgHtmlId(id){
	return `${NORMAL_MSG_TYPE}-${id}`
}
function getSpecialMsgHtmlId(id){
	return `${NORMAL_MSG_TYPE}-special-${id}`
}

function createMsgCard(id, row){
	var icon = `<div id="${getMsgHtmlId(row.ID)}"></div>`
	if(row.IsLiked){icon = `<div id="${getMsgHtmlId(row.ID)}"><a href="#"><i class="fa fa-heart fa-2x" aria-hidden="true"></i></a></div>`}
	var s = `<h3 class="card-title">${row.User}</h3>
					<p class="card-text">${row.Value}</p>
					<p class="card-text" id="${getSpecialMsgHtmlId(row.ID)}">Likes: ${row.LikeNum}</p>
					`;
	s += icon
	return s;
}



function playAnimation(){
	// console.log(elements);
	for(let i=0; i<elements.length; i++){
		$(elements[i]).hide();
		$(elements[i]).fadeIn(i*200);
	}
}

function createCards(rows, cardType, callback){
	var n = rows.length;
	var s = "";
	for(i=0; i<n/3; i++){
		var rowS = "";
		var bound = i*3+3;
		if(i*3+3>n) bound = n;
		for(j=i*3; j<bound ;j++){
			rowS+=createCard(j, rows[j], cardType);
			elements.push('#card-'+cardType+'-'+rows[j].ID);
		}
		s+=closeWithRow(rowS);
	}
	// var preBody = $("#home-table").html();
	$("#home-table").append(s);
	if(callback){
		for(let i=0; i<n/3; i++){
			var bound = i*3+3;
			if(i*3+3>n) bound = n;
			for(let j=i*3; j<bound ;j++){
				$(document).off('click','#card-'+cardType+'-'+rows[j].ID);
				$(document).on('click', '#card-'+cardType+'-'+rows[j].ID, {}, function(e){
					callback(j, rows[j])
				})
			}
		}
	}
}

function normalMsgCardCallback(index, row){
	console.log("normalMsgCardCallback",row)
	var apicall = "/LikeMsg"
	if(row.IsLiked) apicall = "/UnlikeMsg"
	var htmlContent = `<a href="#"><i class="fa fa-heart fa-2x" aria-hidden="true"></i></a>`
	if(row.IsLiked) htmlContent = ""
	// console.log("special html id: ", $("#"+getSpecialMsgHtmlId(row.ID)).html())
	$.post(url+apicall,{msgid:row.ID}).done(function(data){
		console.log("like call success", data)
		console.log(`${getMsgHtmlId(row.ID)}`)
		var cardId = `${getMsgHtmlId(row.ID)}`
		console.log(htmlContent)
		$("#"+cardId).html(htmlContent)
		row.IsLiked = !row.IsLiked
		row.LikeNum--
		if(row.IsLiked) row.LikeNum+=2
		// console.log("likeNum", row.LikeNum)
		$("#"+getSpecialMsgHtmlId(row.ID)).html("Likes: "+row.LikeNum)
		// if(row.IsLiked) row.LikeNum++
		// else row.LikeNum--
	})
}

function musicCardCallback(index, row){
	$.get(url+"/Home/card-",{Aid:row.aid,Tid:row.tid}).done(function(data){
				console.log(data);
 				SingerData=data.SingerData;
				RateData=data.RateData;
				PlaylistData=data.PlaylistData;
				PlaylistNum=PlaylistData.length;
				 $('#CardIink').attr('src' ,row.link);
				 $('#CardTname').html("\""+row.tname+"\"");
				 $('#CardScore').html(""+RateData[0].AveRate+"");
				 $('#CardAname').html(""+SingerData[0].aname+"");
				 for(let i=0; i<PlaylistNum; i++){
				  $('#AddtoPlaylist').html("<a>"+PlaylistData[0].ptitle+"</a>");
			    }
				// $("#PlaySong").html(playWindow);
				 $("#PlaySong").modal('show');
					$(document).off('click','#CardAname');
					$(document).on('click', '#CardAname', {}, function(e){
						$("#PlaySong").modal('hide');;
						artistCallback(0, row);
					})

					 $('#GiveRate a').click(function(){
					   var selText = $(this).text();
					   	 $.get(url+"/Home/rate",{Score:selText,Tid:row.tid}).done(function(data){

                          $('#CardScore').html(""+data.RateData[0].AveRate+"");
					   })

					 });
					 $('#AddtoPlaylist a').click(function(){
						 var selText=$(this).text();
						 $.get(url+"/Home/AddtoPlaylist",{Pid:PlaylistData[0].pid,Tid:row.tid}).done(function(data){

						 })
					 });




			  });
}


function albumCallback(index, row){
	let alid = row.alid;
	var pageAction = function(){
		$.get(url+"/Playlist/GetAlbumSongs",{alid:alid}).done(function(data){
			$("#home-table").html("")
			createCards(data, MUSIC_CARD_TYPE, musicCardCallback);
			playAnimation();
		})
	}
	browseStack.push({action:pageAction, params:{pageTag:"Album"+row.altitle}});
	playStackTop();
}

function playlistCallback(index, row){
	let pid = row.pid;
	var pageAction = function(){
		$.get(url+"/Playlist/GetPlaylistSongs",{pid:pid}).done(function(data){
			$("#home-table").html("");
			createCards(data, MUSIC_CARD_TYPE, musicCardCallback);
			playAnimation();
		})
	}
	browseStack.push({action:pageAction, params:{pageTag:"Playlist"+row.ptitle}});
	playStackTop();
}

var isArtistFollowed;

function setLikeBtnState(isArtistFollowed){
	if(isArtistFollowed==="1"){
		$("#artist-info #btn-like").attr('class','btn btn-secondary');
		$("#artist-info #btn-like").html("Liked");
	}else{
		$("#artist-info #btn-like").attr('class','btn btn-primary');
		$("#artist-info #btn-like").html("Like");
	}
}

function artistCallback(index, row){
	$("#artist-info").modal("show");
	$.get(url+"/Artist/GetArtistInfo",{aid:row.aid}).done(function(data){
		console.log(data);
		let dataRow = data.artistInfoResult[0];
		isArtistFollowed = data.isArtistFollowed;
		$("#artist-info #aname").html("Artist Name: "+dataRow.aname);
		$("#artist-info #adesc").html("Artist Introduction: "+dataRow.adescription);
		setLikeBtnState(isArtistFollowed);

		$(document).off('click',"#artist-info #btn-like");
		$(document).on('click', "#artist-info #btn-like", {}, function(e){
			$.get( url+"/Artist/UpdateLike",{aid:dataRow.aid,newState:isArtistFollowed==="1"?"0":"1"}).done(function( data ) {
				isArtistFollowed = data;
				setLikeBtnState(isArtistFollowed);
		  	});
		})
		$(document).off('click',"#artist-info #btn-tracks");
		$(document).on('click', "#artist-info #btn-tracks", {}, function(e){
			let aid = dataRow.aid;
			var pageAction = function(){
				$.get( url+"/Artist/GetArtistTracks",{aid:aid}).done(function( data ) {
					$('#artist-info').modal('hide');
					$("#home-table").html("");
					createCards(data, MUSIC_CARD_TYPE,musicCardCallback);
					playAnimation();
				});
			}
			browseStack.push({action:pageAction, params:{pageTag: dataRow.aname+" Tracks"}})
			playStackTop();
		})
	})
}


function setFollowBtnState(isUserFollowed){
	if(isUserFollowed==="1"){
		$("#user-info #btn-follow").attr('class','btn btn-secondary');
		$("#user-info #btn-follow").html("Followed");
	}else{
		$("#user-info #btn-follow").attr('class','btn btn-primary');
		$("#user-info #btn-follow").html("Follow");
	}
}

var isUserFollowed;

function userCardCallback(index, row){
	$("#user-info").modal('show');
	$.get(url+"/User/UserInfo",{uname:row.uname}).done(function(data){
		console.log(data);
		let dataRow = data.userInfoData[0];
		isUserFollowed = data.isUserFollowedData;
		$("#user-info #uname").html("User Name: "+dataRow.uname);
		$("#user-info #rname").html("User Real Name: "+dataRow.rname);
		$("#user-info #email").html("User Email: "+dataRow.uemail);
		$("#user-info #city").html("User City: "+dataRow.ucity);
		setFollowBtnState(isUserFollowed);
		$(document).off('click',"#user-info #btn-follow");
		$(document).on('click', "#user-info #btn-follow", {}, function(e){
			$.get( url+"/User/UpdateFollow",{uname:dataRow.uname,newState:isUserFollowed==="1"?"0":"1"}).done(function( data ) {
				isUserFollowed = data;
				setFollowBtnState(isUserFollowed);
		  	});
		})

		$(document).off('click',"#user-info #btn-playlist");
		$(document).on('click', "#user-info #btn-playlist", {}, function(e){
			let userName = dataRow.uname;
			var pageAction = function(){
				$.get( url+"/Playlist/GetUserPlaylist",{userName:userName}).done(function( data ) {
					$('#user-info').modal('hide');
					$("#home-table").html("");
					createCards(data, PLAYLIST_CARD_TYPE,playlistCallback);
					playAnimation();
				});
			}
			browseStack.push({action:pageAction, params:{pageTag: userName+" Playlists"}})
			playStackTop();
		})
	})
}


function getFeeds(){
	$.get( url+"/Home/GetFeeds").done(function( data ) {
		$("#home-table").html("");
 		createCards(data.albumData, NEW_ALBUM_FEED, albumCallback);
 		createCards(data.playData, PLAY_FEED, musicCardCallback);
 		console.log(data);
 		playAnimation();
  	});
}

$(window).on('load', function() {
	getMsgs(LASTEST_MSG_FLAG)
});


$( "#logout" ).click(function() {
  $.get( url+"/User/Logout").done(function( data ) {
 		window.location.replace(url);
  	});
});



$("#search-btn").click(function(){
	let keyword = $("#search-keyword").val();
	var pageAction = function pageSearchAction(){
		$.get(url+"/Home/Search",{keyword:keyword}).done(function(data){
			console.log(data);
			$("#home-table").html("");
			createCards(data.trackData, MUSIC_CARD_TYPE,musicCardCallback);
			createCards(data.artistData, ARTIST_CARD_TYPE, artistCallback);
			createCards(data.albumData, ALBUM_CARD_TYPE,albumCallback);
			createCards(data.playlistData, PLAYLIST_CARD_TYPE,playlistCallback);
			createCards(data.userData, USER_CARD_TYPE,userCardCallback);
			playAnimation();
		})
	}
	browseStack.push({action:pageAction, params:{pageTag:"Search "+"\""+keyword+"\""}});
	playStackTop();
});

$( "#profile" ).click(function() {
  $.get( url+"/Home/Profile").done(function( data ) {
	  $('#Name').attr('placeholder' ,data[0].rname);
	  $('#Email').attr('placeholder' ,data[0].uemail);
	  $('#City').attr('placeholder' ,data[0].ucity);
	});
});

$( "#UpdateUserProfile" ).click(function() {
  $.get( url+"/Home/UpdateUserProfile",
  {realName:$("#Name").val() ,email:$("#Email").val(),city:$("#City").val()})
  .done(function( data ){})
});

$("#showPlaylist").click(function(){
	var pageAction = function(){
		$.get( url+"/Playlist/GetUserPlaylist",{userName:local_data}).done(function( data ) {
			$("#home-table").html("");
			createCards(data, PLAYLIST_CARD_TYPE,playlistCallback);
			playAnimation();
		});
	}
	browseStack.push({action:pageAction, params:{pageTag:"Your Playlists"}})
	playStackTop();
})

$("#yourFollows").click(function(){
	var pageAction = function(){
		$.get( url+"/UserBehavior/GetFollows").done(function( data ) {
			$("#home-table").html("");
			createCards(data, USER_CARD_TYPE,userCardCallback);
			playAnimation();
		});
	}
	browseStack.push({action:pageAction, params:{pageTag:"Your Follows"}});
	playStackTop();
})

$("#yourArtists").click(function(){
	var pageAction = function(){
		$.get( url+"/Artist/ShowUserArtists",{userName:local_data}).done(function( data ) {
			$("#home-table").html("");
			createCards(data, ARTIST_CARD_TYPE,artistCallback);
			playAnimation();
		});
	}
	browseStack.push({action:pageAction, params:{pageTag:"Your Artists"}});
	playStackTop();
})


$("#btn-add-playlist").click(function(){
	$("#custom-playlist").modal("hide");
	$.get( url+"/Playlist/AddPlaylist",{ptitle:$("#ptitle").val()}).done(function( data ) {

	});
})

$("#back-btn").click(function(){
	$("#msg-board").modal("show");
})

$("#post-msg-btn").click(function(){
	console.log("post-msg-btn", " pressed")
	$.post( url+"/SendMsg",{value:$("#msg-board-input").val()}).done(function( data ) {
			console.log(data);
			if(data=="0"){
				alert("post msg fail, try again")
			}else{
				$("#msg-board").modal("hide");
				refreshTable()
			}
		});
})

$("#post-user-del").click(function(){
	$.post( url+"/DelUser",{}).done(function( data ) {
			console.log(data);
			if(data == "0"){
				alert("Delete Your Account Successfully")
				window.location.replace(url + "/index.html");
			}else{

			}
		});
})

$("#get-more-msg-btn").click(function(){
	if(lastMsgId==0){return}
	getMsgs(lastMsgId-1)
})


function getMsgs(msgId){
	console.log("msgId",msgId)
	$.get( url+"/GetMsg",{index:msgId}).done(function( data ) {
			console.log(data);
			var parsedData = JSON.parse(data)
			if(parsedData!=null && parsedData.length>0){
				var tempLastMsgId = parsedData[parsedData.length-1].ID
				lastMsgId = tempLastMsgId
				createCards(parsedData, NORMAL_MSG_TYPE, normalMsgCardCallback)
			}
		});
}



$("#get-like-url").click(function(){
		window.location.replace(url + "/view/like.html");
});


$("#refresh-btn").click(function(){
	refreshTable()
})

function refreshTable(){
	$("#home-table").html("");
	lastMsgId = 0
	getMsgs(LASTEST_MSG_FLAG)
}

$("#user-follow-search-bar-btn").click(function(){
	console.log("follow btn press")
	$.post( url+"/FollowUser",{}).done(function( data ) {
			console.log(data)
	});
})
