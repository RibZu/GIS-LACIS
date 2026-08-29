function togglePasswordVisibility(button) {
    var input = document.getElementById(button.dataset.target);
    if (!input) return;
    var showing = input.type === 'text';
    input.type = showing ? 'password' : 'text';
    var icon = button.querySelector('i');
    if (icon) {
        icon.classList.toggle('bi-eye', showing);
        icon.classList.toggle('bi-eye-slash', !showing);
    }
}

document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('[data-toggle-password]').forEach(function (btn) {
        btn.addEventListener('click', function () {
            togglePasswordVisibility(btn);
        });
    });

    var pass = document.getElementById('password');
    var repeat = document.getElementById('password_confirm');
    if (pass && repeat) {
        var validateMatch = function () {
            repeat.setCustomValidity(repeat.value !== pass.value ? 'Las contraseñas no coinciden' : '');
        };
        pass.addEventListener('input', validateMatch);
        repeat.addEventListener('input', validateMatch);
    }
});
