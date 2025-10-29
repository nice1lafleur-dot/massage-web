(function(){
  const btn = document.querySelector('.nav-toggle');
  const nav = document.querySelector('nav.nav');
  if(!btn || !nav) return;
  btn.addEventListener('click', function(){
    nav.classList.toggle('open');
    const expanded = nav.classList.contains('open');
    btn.setAttribute('aria-expanded', expanded ? 'true' : 'false');
  });
  // close when clicking outside
  document.addEventListener('click', function(e){
    if(!nav.classList.contains('open')) return;
    if(nav.contains(e.target) || btn.contains(e.target)) return;
    nav.classList.remove('open');
    btn.setAttribute('aria-expanded','false');
  });
})();