document.addEventListener('DOMContentLoaded', function () {
  window.addEventListener('scroll', function (e) {
    const scrollY = window.scrollY;
    const nav = document.getElementById('nav');
    const navTop = nav.offsetTop;
    const toolbar = document.getElementById('toolbar');
    console.log(toolbar.offsetTop);
    const toolbarTop = toolbar.offsetTop - 50;
    const posts = document.getElementById('posts');
    const postsTop = posts.offsetTop;
    const postsWidth = posts.offsetWidth - 10;
    const tclass = ['fixed', 'top-3', 'z-10'];
    if (scrollY < postsTop) {
      tclass.forEach((c) => {
        toolbar.classList.remove(c);
      });
      toolbar.style.width = '';
      // nav.classList.remove('bg-white');
    } else {
      tclass.forEach((c) => {
        toolbar.classList.add(c);
      });
      toolbar.style.width = `${postsWidth}px`;
      // nav.classList.add('bg-white');
    }
  });
});
