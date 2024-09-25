document.addEventListener('DOMContentLoaded', function () {
  document.getElementById('logout').addEventListener('click', function () {
    Cookies.remove('auth-session');
    window.location.href = '/auth/logout';
  });
});
