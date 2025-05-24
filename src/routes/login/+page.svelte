<script>
  import { goto } from '$app/navigation';
  
  let email = '';
  let password = '';
  let rememberMe = false;
  let errorMessage = '';
  
  function handleLogin() {
    // For demonstration purposes only
    // In a real app, you would validate credentials with a backend
    if (email && password) {
      // Simulate successful login
      goto('/'); // Redirect to home page after login
    } else {
      errorMessage = 'Please enter both email and password';
    }
  }
</script>

<svelte:head>
  <title>Login - LynxLinkage</title>
  <meta name="description" content="Log in to your LynxLinkage account" />
</svelte:head>

<div class="login-page">
  <header>
    <div class="logo">
      <a href="/">LynxLinkage</a>
    </div>
  </header>
  
  <main>
    <div class="login-container">
      <h1>Welcome Back</h1>
      <p class="subtitle">Log in to your account to continue</p>
      
      {#if errorMessage}
        <div class="error-message">{errorMessage}</div>
      {/if}
      
      <form on:submit|preventDefault={handleLogin}>
        <div class="form-group">
          <label for="email">Email</label>
          <input 
            type="email" 
            id="email" 
            bind:value={email} 
            placeholder="Enter your email"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            type="password" 
            id="password" 
            bind:value={password} 
            placeholder="Enter your password"
            required
          />
        </div>
        
        <div class="form-options">
          <label class="remember-me">
            <input type="checkbox" bind:checked={rememberMe} />
            <span>Remember me</span>
          </label>
          
          <a href="/forgot-password" class="forgot-password">Forgot password?</a>
        </div>
        
        <button type="submit" class="login-button">Log In</button>
      </form>
      
      <div class="signup-prompt">
        <p>Don't have an account? <a href="/signup">Sign up</a></p>
      </div>
    </div>
  </main>
  
  <footer>
    <p>&copy; {new Date().getFullYear()} LynxLinkage. All rights reserved.</p>
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
  
  /* Login Page Layout */
  .login-page {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }
  
  /* Header Styles */
  header {
    padding: 1.5rem 5%;
    background-color: rgba(255, 255, 255, 0.95);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
  }
  
  .logo {
    font-size: 1.5rem;
    font-weight: 700;
  }
  
  .logo a {
    text-decoration: none;
    color: #6366f1;
  }
  
  /* Main Content Styles */
  main {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 2rem 1rem;
  }
  
  .login-container {
    background-color: white;
    border-radius: 10px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);
    padding: 2.5rem;
    width: 100%;
    max-width: 450px;
  }
  
  h1 {
    font-size: 1.875rem;
    color: #1a202c;
    margin-bottom: 0.5rem;
    text-align: center;
  }
  
  .subtitle {
    color: #6b7280;
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .error-message {
    background-color: #fee2e2;
    color: #b91c1c;
    padding: 0.75rem;
    border-radius: 0.375rem;
    margin-bottom: 1.5rem;
    text-align: center;
  }
  
  /* Form Styles */
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #4b5563;
  }
  
  input[type="email"],
  input[type="password"] {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.375rem;
    background-color: #f9fafb;
    transition: border-color 0.2s;
    font-size: 1rem;
  }
  
  input[type="email"]:focus,
  input[type="password"]:focus {
    outline: none;
    border-color: #6366f1;
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
  }
  
  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    font-size: 0.875rem;
  }
  
  .remember-me {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
  
  .forgot-password {
    color: #6366f1;
    text-decoration: none;
  }
  
  .forgot-password:hover {
    text-decoration: underline;
  }
  
  .login-button {
    width: 100%;
    padding: 0.75rem 1.5rem;
    background-color: #6366f1;
    color: white;
    border: none;
    border-radius: 0.375rem;
    font-weight: 500;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s ease, transform 0.2s ease;
  }
  
  .login-button:hover {
    background-color: #4f46e5;
    transform: translateY(-2px);
  }
  
  .login-button:active {
    transform: translateY(0);
  }
  
  /* Sign Up Prompt */
  .signup-prompt {
    margin-top: 2rem;
    text-align: center;
    font-size: 0.875rem;
    color: #6b7280;
  }
  
  .signup-prompt a {
    color: #6366f1;
    text-decoration: none;
    font-weight: 500;
  }
  
  .signup-prompt a:hover {
    text-decoration: underline;
  }
  
  /* Footer Styles */
  footer {
    text-align: center;
    padding: 1.5rem;
    color: #9ca3af;
    font-size: 0.875rem;
    background-color: white;
    border-top: 1px solid #f3f4f6;
  }
  
  /* Responsive Styles */
  @media (max-width: 640px) {
    .login-container {
      padding: 2rem 1.5rem;
    }
    
    h1 {
      font-size: 1.5rem;
    }
  }
</style>
