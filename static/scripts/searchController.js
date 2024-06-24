const searchForm = document.querySelector(".searchForm");
let searchBarVisible = false;

function searchHandler() {
	if (searchBarVisible == false){
		searchForm.style.display = "flex";
		searchBarVisible = true;
	} else {
		searchForm.style.display = "none";
		searchBarVisible = false;
	}
} 

searchForm.addEventListener("submit", async function(event){
	event.preventDefault();
	const formData = new FormData(this)
	const searchBox = document.querySelector(".searchBox")
	const jsonObject = {};
	const endpoint = "/search";

	if (searchBox.value == '') {
		fetchDataAndRender();
		return;
	}

	jsonObject["id"] = null;
	jsonObject["quantity"] = null;
	jsonObject["name"] = searchBox.value;
	jsonObject["category"] = searchBox.value;

	await fetch(endpoint, {
  		method: 'POST',
  		body: JSON.stringify(jsonObject)
 	}).then( async (data) => {
 		const jsonResponse = await data.json();
 		render(jsonResponse);
 	})

});