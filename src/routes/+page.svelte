<script>
  // Images for gallery with links
  const galleryImages = [
    { 
      src: 'https://placehold.co/600x400?text=Product+View+1', 
      alt: 'Product view 1',
      url: 'https://example.com/feature1'
    },
    { 
      src: 'https://placehold.co/600x400?text=Product+View+2', 
      alt: 'Product view 2',
      url: 'https://example.com/feature2'
    },
    { 
      src: 'https://placehold.co/600x400?text=Product+View+3', 
      alt: 'Product view 3',
      url: 'https://example.com/feature3'
    },
    { 
      src: 'https://placehold.co/600x400?text=Product+View+4', 
      alt: 'Product view 4',
      url: 'https://example.com/feature4'
    },
  ];
</script>

<svelte:head>
  <title>LynxLinkage - Connect Effortlessly</title>
  <meta name="description" content="A seamless way to connect and collaborate" />
</svelte:head>

<div class="landing-page">
  <!-- Header -->
  <header>
    <div class="logo">LynxLinkage</div>
    <nav>
      <ul>
        <li><a href="#gallery">Gallery</a></li>
        <li><a href="#about">About</a></li>
        <li><a href="#contact" class="contact-btn">Contact</a></li>
        <li><a href="/login" class="login-btn">Login</a></li>
      </ul>
    </nav>
  </header>

  <!-- Hero Section -->
  <section class="hero">
    <div class="hero-content">
      <h1>Effortless Connections, Endless Possibilities</h1>
      <p>Experience a new way of connecting with your world. Simple, intuitive, and designed for you.</p>
      <div class="cta-buttons">
        <button class="primary-btn">Get Started</button>
        <button class="secondary-btn">Learn More</button>
      </div>
    </div>
    <div class="hero-image">
      <img src="https://placehold.co/800x600?text=Hero+Image" alt="LynxLinkage product showcase" />
    </div>
  </section>

  <!-- Gallery Section -->
  <section id="gallery" class="gallery">
    <h2>Discover Our Features</h2>
    <p>Take a closer look at what makes LynxLinkage special.</p>
    
    <div class="gallery-grid">
      {#each galleryImages as image}
        <div class="gallery-item">
          <a href={image.url} target="_blank" rel="noopener noreferrer">
            <img src={image.src} alt={image.alt} />
          </a>
        </div>
      {/each}
    </div>
  </section>

  <!-- Footer -->
  <footer>
    <div class="footer-content">
      <div class="footer-logo">
        <div class="logo">LynxLinkage</div>
        <p>Connecting the world, one link at a time.</p>
      </div>
      
      <div class="footer-links">
        <div class="link-group">
          <h3>Product</h3>
          <ul>
            <li><a href="#pricing">Pricing</a></li>
            <li><a href="#faq">FAQ</a></li>
          </ul>
        </div>
        
        <div class="link-group">
          <h3>Company</h3>
          <ul>
            <li><a href="#about">About Us</a></li>
            <li><a href="#contact">Contact</a></li>
          </ul>
        </div>
      </div>
    </div>
    
    <div class="copyright">
      <p>&copy; {new Date().getFullYear()} LynxLinkage. All rights reserved.</p>
    </div>
  </footer>
</div>

<style>
  /* Global Styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    color: #333;
    background-color: #fafafa;
  }
  
  .landing-page {
    max-width: 100%;
    overflow-x: hidden;
  }
  
  h1, h2, h3 {
    font-weight: 600;
    color: #2d3748;
  }
  
  /* Header Styles */
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 5%;
    background-color: rgba(255, 255, 255, 0.95);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
    position: sticky;
    top: 0;
    z-index: 100;
  }
  
  .logo {
    font-size: 1.5rem;
    font-weight: 700;
    color: #6366f1;
  }
  
  nav ul {
    display: flex;
    list-style: none;
    gap: 2rem;
    margin: 0;
    padding: 0;
  }
  
  nav a {
    text-decoration: none;
    color: #4b5563;
    font-weight: 500;
    transition: color 0.3s ease;
  }
  
  nav a:hover {
    color: #6366f1;
  }
  
  .contact-btn {
    padding: 0.5rem 1rem;
    background-color: #f3f4f6;
    border-radius: 0.375rem;
  }
  
  .contact-btn:hover {
    background-color: #e5e7eb;
  }
  
  .login-btn {
    padding: 0.5rem 1rem;
    background-color: #6366f1;
    color: white;
    border-radius: 0.375rem;
    transition: background-color 0.3s ease;
  }
  
  .login-btn:hover {
    background-color: #4f46e5;
  }
  
  /* Hero Section Styles */
  .hero {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 5rem 5%;
    gap: 2rem;
  }
  
  .hero-content {
    flex: 1;
    max-width: 600px;
  }
  
  .hero h1 {
    font-size: 3rem;
    line-height: 1.2;
    margin-bottom: 1.5rem;
    color: #1a202c;
  }
  
  .hero p {
    font-size: 1.125rem;
    line-height: 1.6;
    color: #4b5563;
    margin-bottom: 2rem;
  }
  
  .cta-buttons {
    display: flex;
    gap: 1rem;
  }
  
  .primary-btn, .secondary-btn {
    padding: 0.75rem 1.5rem;
    border-radius: 0.375rem;
    font-weight: 500;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.3s ease;
    border: none;
  }
  
  .primary-btn {
    background-color: #6366f1;
    color: white;
  }
  
  .primary-btn:hover {
    background-color: #4f46e5;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
  }
  
  .secondary-btn {
    background-color: #f3f4f6;
    color: #4b5563;
  }
  
  .secondary-btn:hover {
    background-color: #e5e7eb;
    transform: translateY(-2px);
  }
  
  .hero-image {
    flex: 1;
    display: flex;
    justify-content: center;
  }
  
  .hero-image img {
    max-width: 100%;
    height: auto;
    border-radius: 10px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  }
  
  /* Gallery Section Styles */
  .gallery {
    padding: 5rem 5%;
    background-color: #f9fafb;
    text-align: center;
  }
  
  .gallery h2 {
    font-size: 2.25rem;
    margin-bottom: 1rem;
  }
  
  .gallery p {
    font-size: 1.125rem;
    color: #6b7280;
    max-width: 600px;
    margin: 0 auto 3rem;
  }
  
  .gallery-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 2rem;
  }
  
  .gallery-item {
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    cursor: pointer;
  }
  
  .gallery-item:hover {
    transform: translateY(-5px);
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  }
  
  .gallery-item a {
    display: block;
    width: 100%;
    height: 100%;
  }
  
  .gallery-item img {
    width: 100%;
    height: auto;
    display: block;
  }
  
  /* Footer Styles */
  footer {
    background-color: #f3f4f6;
    padding: 4rem 5% 2rem;
    color: #4b5563;
  }
  
  .footer-content {
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 3rem;
    margin-bottom: 3rem;
  }
  
  .footer-logo {
    flex: 1;
    min-width: 250px;
  }
  
  .footer-logo p {
    margin-top: 1rem;
    font-size: 0.875rem;
    color: #6b7280;
  }
  
  .footer-links {
    display: flex;
    flex-wrap: wrap;
    gap: 3rem;
  }
  
  .link-group h3 {
    font-size: 1rem;
    margin-bottom: 1.25rem;
    color: #374151;
  }
  
  .link-group ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .link-group li {
    margin-bottom: 0.75rem;
  }
  
  .link-group a {
    text-decoration: none;
    color: #6b7280;
    font-size: 0.875rem;
    transition: color 0.3s ease;
  }
  
  .link-group a:hover {
    color: #6366f1;
  }
  
  .copyright {
    text-align: center;
    border-top: 1px solid #e5e7eb;
    padding-top: 2rem;
    font-size: 0.875rem;
    color: #9ca3af;
  }
  
  /* Responsive Styles */
  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
      text-align: center;
      padding: 3rem 5%;
    }
    
    .hero-content {
      max-width: 100%;
    }
    
    .cta-buttons {
      justify-content: center;
    }
    
    .hero-image {
      margin-top: 2rem;
    }
    
    .footer-content {
      flex-direction: column;
      gap: 2rem;
    }
    
    nav ul {
      gap: 1rem;
    }
    
    header {
      padding: 1rem 5%;
      flex-direction: column;
      gap: 1rem;
    }
  }
</style>
