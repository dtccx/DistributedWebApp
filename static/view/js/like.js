var url = "http://localhost:8080";
var lastMsgId = 0

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
	// var preBody = $("#home-table").html();
	$("#home-table").append(s);
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
					<div class="card-block home-table-card" id="card-${cardType}-${id}">${ret}</div>
				</div>
			 </div>`;
	return temp;
}

function createMsgCard(id, row){
	var s = `<h3 class="card-title">${row.User}</h3>
					<p class="card-text">${row.Value}</p>
					`;
	return s;
}


$("#get-like-btn").click(function(){
	if(lastMsgId == 0){return}
  getMsgs(lastMsgId - 1)
})

function getLikeMsgs(msgId){
	console.log("msgId",msgId)

	$.get( url+"/LikeList",{}).done(function( data ) {
			console.log(data);
			var parsedData = JSON.parse(data)
			if(parsedData != null && parsedData.length > 0){
				var tempLikeMsgId = parsedData[parsedData.length-1].ID
				lastMsgId = tempLikeMsgId
				createCards(parsedData, NORMAL_MSG_TYPE, null)
			}
		});
}
