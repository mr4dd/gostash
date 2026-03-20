const addButtonPressed = document.querySelector(".new");
const InfoContainer = document.querySelector(".infoContainer");
const form = document.querySelector(".formData");
const entries = document.querySelector(".entries");
let isVisible = false;

document.body.addEventListener('keydown', function(e) {
  if (e.key == "Escape") {
  	isVisible = true;
  	ContainerDisplayManager();
  }
});

function ContainerDisplayManager() {
	if (isVisible){
		InfoContainer.style.display = "none";
		isVisible = false;
	}else{
		InfoContainer.style.display = "flex";
		isVisible = true;
	}	
};

function newEntry(){
	form.children[5].value = 0;
	form.children[0].value = parseInt(entries.children[entries.children.length - 1].getAttribute("data-id")) + 1 || 0;
	ContainerDisplayManager();
}

function clearForm() {
	const children = form.children;
	for (let i = 0; i < children.length - 1; i++){
		children[i].value = '';
	}
}

form.addEventListener("submit", function(event) {
  let endpoint = '';
  event.preventDefault(); // Prevent form submission
  // Get form data
  const formData = new FormData(this);
  
  // Convert form data to JSON object
  const jsonObject = {};
  formData.forEach((value, key) => {
  	if (key == "newselector") {
 		if (value == '0') {
    	  	endpoint = '/additem';
    	  } else {
   	  	endpoint = '/edititem';
    	  }
	}else if (key == "id" || key == "quantity") {
  		jsonObject[key] = parseInt(value);	
  	}
  	else {
    	jsonObject[key] = value;
   }
  });
  
  fetch(endpoint, {
  	method: 'POST',
  	body: JSON.stringify(jsonObject)
  })
  clearForm()
  ContainerDisplayManager();
  fetchDataAndRender();
});

const editButtonPressed = (button) => {
    const form = document.querySelector(".formData");
    
    const card = button.closest('.entrycontainer');

    const idValue          = card.getAttribute("data-id");
    const nameValue        = card.querySelector('.js-name').innerText;
    const categoryValue    = card.querySelector('.js-category').innerText;
    const quantityValue    = card.querySelector('.js-quantity').innerText;
    const descriptionValue = card.querySelector('.js-description').innerText;

    form.children[0].value = idValue;
    form.children[1].value = quantityValue;
    form.children[2].value = nameValue;
    
    const categoryField = form.children[3].disabled ? form.children[4] : form.children[3];
    categoryField.value = categoryValue;

    form.children[4].value = 1; 
    form.children[5].value = 1;
    form.children[6].value = descriptionValue;

    ContainerDisplayManager();
};

function remove(button) {
    if (!confirm("Are you sure you want to delete this?")) return;

    const card = button.closest('.entrycontainer');
    const id = card.getAttribute("data-id");

    fetch("/deleteitem", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: parseInt(id), quantity: 0, name: "null", category: "null" })
    }).then(() => card.remove());
}
