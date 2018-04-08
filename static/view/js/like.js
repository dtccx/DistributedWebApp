var url = "http://localhost:8080";

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
	if(lastMsgId==0){return}
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
