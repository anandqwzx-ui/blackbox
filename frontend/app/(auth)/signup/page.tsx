'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { FormEvent, useState } from 'react';

import { signup } from '@/lib/api';

export default function SignupPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const response = await signup(email.trim(), password);
      localStorage.setItem('mockapi_token', response.token);
      localStorage.setItem('mockapi_email', response.email);
      router.push('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unable to sign up');
    } finally {
      setLoading(false);
    }
  };

  return (
    <section style={{ maxWidth: '420px', margin: '0 auto' }}>
      <h1>Create an account</h1>
      <p className="lead">Join Mock API Studio to build, test, and share API simulations effortlessly.</p>

      <form onSubmit={handleSubmit} className="card">
        <label htmlFor="email">Email address</label>
        <input
          id="email"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
        />

        <label htmlFor="password">Password</label>
        <input
          id="password"
          type="password"
          placeholder="At least 8 characters"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          minLength={8}
          required
        />

        {error && <p className="error-text">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Creating account…' : 'Sign up'}
        </button>
      </form>

      <p className="muted" style={{ marginTop: '1.5rem' }}>
        Already registered?{' '}
        <Link href="/login" style={{ textDecoration: 'underline' }}>
          Log in
        </Link>
      </p>
    </section>
  );
}
