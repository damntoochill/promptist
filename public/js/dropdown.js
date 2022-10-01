document.addEventListener('mouseup', function(e) {
    
    var el = document.getElementById("userDropdown");
    el.style.display = 'none';
    var el = document.getElementById("browseDropdown");
    el.style.display = 'none';
  });

function userButton() {  
    let el = document.getElementById("userDropdown");
    if (el.style.display == 'none' || el.style.display == '') {
        el.style.display = 'block';
    }else{
        el.style.display = 'none';
    }
}  

function browseButton() {  
    let el = document.getElementById("browseDropdown");
    if (el.style.display == 'none' || el.style.display == '') {
        el.style.display = 'block';
    }else{
        el.style.display = 'none';
    }
}  


