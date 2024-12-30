// form.js
document.getElementById('show-signup').addEventListener('click', function (e) {
    e.preventDefault(); // Prevent default anchor behavior
    document.getElementById('signin').classList.remove('active');
    document.getElementById('signup').style.display = 'block'; // Show sign-up form
    document.querySelector('.form-wrapper').style.transform = 'translateX(-100%)'; // Slide left
});

document.getElementById('show-signin').addEventListener('click', function (e) {
    e.preventDefault(); // Prevent default anchor behavior
    document.getElementById('signup').style.display = 'none'; // Hide sign-up form
    document.getElementById('signin').classList.add('active');
    document.querySelector('.form-wrapper').style.transform = 'translateX(0)'; // Slide back to original position
});