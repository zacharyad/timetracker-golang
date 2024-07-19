document.addEventListener('DOMContentLoaded', () => {
  const menuToggle = document.querySelector('.menu-toggle');
  const navLinks = document.querySelector('.nav-links');
  const modal = document.getElementById('appointmentModal');
  const bookButtons = document.querySelectorAll(
    '.book-appointment, .cta-button'
  );
  const closeModal = document.querySelector('.close-modal');
  const appointmentForm = document.getElementById('appointmentForm');

  menuToggle.addEventListener('click', () => {
    navLinks.classList.toggle('active');
    menuToggle.classList.toggle('active');
    menuToggle.setAttribute(
      'aria-expanded',
      menuToggle.classList.contains('active')
    );
  });

  bookButtons.forEach((button) => {
    button.addEventListener('click', () => {
      modal.style.display = 'block';
      document.body.style.overflow = 'hidden';
    });
  });

  closeModal.addEventListener('click', () => {
    modal.style.display = 'none';
    document.body.style.overflow = 'visible';
  });

  window.addEventListener('click', (e) => {
    if (e.target === modal) {
      modal.style.display = 'none';
      document.body.style.overflow = 'visible';
    }
  });

  appointmentForm.addEventListener('submit', (e) => {
    e.preventDefault();

    const formData = new FormData(appointmentForm);
    console.log('Appointment booked:', Object.fromEntries(formData));
    alert(
      'Thank you for booking an appointment! We will contact you soon to confirm.'
    );
    modal.style.display = 'none';
    document.body.style.overflow = 'visible';
    appointmentForm.reset();
  });

  document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
    anchor.addEventListener('click', function (e) {
      e.preventDefault();
      document.querySelector(this.getAttribute('href')).scrollIntoView({
        behavior: 'smooth',
      });
    });
  });

  const faders = document.querySelectorAll('.fade-in');
  const appearOptions = {
    threshold: 0.5,
    rootMargin: '0px 0px -100px 0px',
  };

  const appearOnScroll = new IntersectionObserver((entries, appearOnScroll) => {
    entries.forEach((entry) => {
      if (!entry.isIntersecting) return;
      entry.target.classList.add('appear');
      appearOnScroll.unobserve(entry.target);
    });
  }, appearOptions);

  faders.forEach((fader) => {
    appearOnScroll.observe(fader);
  });

  gsap.from('.hero h1', { opacity: 0, y: 50, duration: 1, delay: 0.5 });
  gsap.from('.hero p', { opacity: 0, y: 50, duration: 1, delay: 0.8 });
  gsap.from('.cta-button', { opacity: 0, y: 50, duration: 1, delay: 1.1 });

  gsap.from('.service-card', {
    opacity: 0,
    y: 50,
    duration: 0.8,
    stagger: 0.2,
    scrollTrigger: {
      trigger: '.services',
      start: 'top 80%',
    },
  });

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      if (modal.style.display === 'block') {
        modal.style.display = 'none';
        document.body.style.overflow = 'visible';
      }
      if (navLinks.classList.contains('active')) {
        navLinks.classList.remove('active');
        menuToggle.classList.remove('active');
        menuToggle.setAttribute('aria-expanded', 'false');
      }
    }
  });

  document
    .querySelectorAll('.service-card, .about-content, .contact-content')
    .forEach((el) => {
      el.classList.add('fade-in');
    });

  if ('loading' in HTMLImageElement.prototype) {
    const images = document.querySelectorAll('img[loading="lazy"]');
    images.forEach((img) => {
      img.src = img.dataset.src;
    });
  } else {
    const script = document.createElement('script');
    script.src =
      'https://cdnjs.cloudflare.com/ajax/libs/lazysizes/5.3.2/lazysizes.min.js';
    document.body.appendChild(script);
  }

  const inputs = document.querySelectorAll('input, textarea');
  inputs.forEach((input) => {
    input.addEventListener('invalid', (e) => {
      e.preventDefault();
      input.classList.add('error');
    });
    input.addEventListener('input', () => {
      input.classList.remove('error');
    });
  });
});
