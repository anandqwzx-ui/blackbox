import Link from 'next/link';

export default function HomePage() {
  return (
    <>
      <header className="space-y-4">
        <h1>Crudbox</h1>
        <p className="lead">
          Spin up realistic API prototypes in minutes. Design endpoints, manage organisations and projects, and instantly respond to incoming requests.
        </p>
        <div className="flex gap-4 flex-wrap">
          <Link href="/signup" className="inline-flex items-center">
            <button>Sign up</button>
          </Link>
          <Link href="/login" className="inline-flex items-center">
            <button className="bg-transparent border border-white text-white">Log in</button>
          </Link>
        </div>
      </header>

      <section className="mt-16">
        <h2>How it works</h2>
        <div className="grid">
          <div className="card">
            <h3>1. Join</h3>
            <p>Create an account with your email and a secure password.</p>
          </div>
          <div className="card">
            <h3>2. Organise</h3>
            <p>Group your work by organisations and projects with shareable codes.</p>
          </div>
          <div className="card">
            <h3>3. Mock</h3>
            <p>Define HTTP endpoints, customise payloads, and instantly serve mocked responses.</p>
          </div>
        </div>
      </section>
    </>
  );
}
