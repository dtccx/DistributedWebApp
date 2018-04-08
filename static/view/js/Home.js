var url = "http://localhost:8080";

var browseStack = []
var elements = []

var MUSIC_CARD_TYPE = 0;
var ARTIST_CARD_TYPE = 1;
var ALBUM_CARD_TYPE = 2;
var PLAYLIST_CARD_TYPE = 3;
var USER_CARD_TYPE = 4;
var NEW_ALBUM_FEED = 5;
var PLAY_FEED = 6;

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
		case MUSIC_CARD_TYPE:
			cardClass = "track-card";
			ret = createMusicCard(id, row);
			break;
		case ARTIST_CARD_TYPE:
			cardClass = "artist-card";
			ret = createArtistCard(id, row);
			break;
		case ALBUM_CARD_TYPE:
			cardClass = "album-card";
			ret = createAlbumCard(id, row);
			break;
		case PLAYLIST_CARD_TYPE:
			cardClass = "playlist-card";
			ret = createPlaylistCard(id, row);
			break;
		case USER_CARD_TYPE:
			cardClass = "artist-card";
			ret = createUserCard(id, row);
			break;
		case NEW_ALBUM_FEED:
			cardClass = "artist-card";
			ret = createNewAlbumFeed(id, row);
			break;
		case PLAY_FEED:
			cardClass = "track-card";
			ret = createPlayFeed(id, row);
			break;
		default:
			break;
	}

	var temp = `<div class="col-sm-4">
				<div class="card home-table-card2 ${cardClass}">
					<div class="card-block home-table-card" id="card-${cardType}-${id}">${ret}</div>
				</div>
			 </div>`;
	return temp;
}

function getTimeInterval(seconds){
	var timeInterval = 0; var unit = "days";
	if(seconds>(365*24*3600)){
		timeInterval = seconds/(365*24*3600);
		unit = "years";
	}
	else if(seconds>(30*24*3600)){
		timeInterval = seconds/(30*24*3600);
		unit = "months";
	}
	else{
		timeInterval = seconds/(24*3600);
		unit = "days"
	}
	return Math.floor(timeInterval) +" "+ unit;
}

function createPlayFeed(id, row){
	var timeInterval = getTimeInterval(parseInt(row.timeinterval));
	var s = `<i class="fa fa-commenting-o fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.uname} just listen to "${row.tname}" ${timeInterval} ago!</h3>`;
	return s;
}

function createNewAlbumFeed(id,row){
	// console.log(row);
	var timeInterval = getTimeInterval(parseInt(row.timeinterval));
	var s = `<i class="fa fa-list-alt fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.aname} just release the new album "${row.altitle}" ${timeInterval} ago!</h3>`;
    return s;
}

function createPlaylistCard(id, row){
	var s = `<i class="fa fa-list-alt fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.ptitle}</h3>
					<p class="card-text">Owned By:${row.uname}</p>
					<p class="card-text">Release Time:${row.buildtime}</p>`;
	return s;
}

function createAlbumCard(id, row){
	var s = `<i class="fa fa-list-alt fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.altitle}</h3>
					<p class="card-text">Release Time:${row.Albumtime}</p>`;
	return s;
}

function createMusicCard(id,row){
	// console.log(row);
	var s = `<i class="fa fa-music fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.tname}</h3>
					<p class="card-text">By:${row.aname}</p>
					<p class="card-text">Genre:${row.tgenre}</p>
					<p class="card-text">Duration:${row.duration} seconds</p>`;
    return s;
}

function createArtistCard(id, row){
	var s = `<i class="fa fa-user-circle fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.aname}</h3>
					<p class="card-text">Description:${row.adescription}</p>`;
    return s;
}

function createUserCard(id, row){
	var s = `<i class="fa fa-user-o fa-2x" aria-hidden="true"></i>
					<h3 class="card-title">${row.uname}</h3>
					<p class="card-text">Name:${row.rname}</p>
					<p class="card-text">City:${row.ucity}</p>`;
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
			elements.push('#card-'+cardType+'-'+j);
		}
		s+=closeWithRow(rowS);
	}
	var preBody = $("#home-table").html();
	$("#home-table").html(preBody+s);
	if(callback){
		for(let i=0; i<n/3; i++){
			var bound = i*3+3;
			if(i*3+3>n) bound = n;
			for(let j=i*3; j<bound ;j++){
				$(document).off('click','#card-'+cardType+'-'+j);
				$(document).on('click', '#card-'+cardType+'-'+j, {}, function(e){
					callback(j, rows[j])
				})
			}
		}
	}
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
	document.addEventListener('keyup', (event) => {
	  		const keyName = event.key;
	  		if (keyName === 'a') {
	  			console.log(browseStack);
  		}
	}, false);

	// $("#Greeting").html("Welcome! "+local_data);
	playStackTop();
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
	$("#btnLogin").click(function(){
		$.post( url+"/SendMsg",{value:$("#msg-board-input").val()}).done(function( data ) {
	      console.log(data);
				if(data=="0"){
					alert("post msg fail, try again")
				}else{
					$("#msg-board").modal("hide");
				}
	  	});
	});
})
