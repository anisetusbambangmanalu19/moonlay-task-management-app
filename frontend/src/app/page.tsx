import { redirect } from 'next/navigation';

// Root page redirects to /tasks (which will redirect to /login if not authenticated)
export default function Home() {
  redirect('/tasks');
}
