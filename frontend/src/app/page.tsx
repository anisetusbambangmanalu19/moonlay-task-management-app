import { redirect } from 'next/navigation';

// Halaman root mengarahkan ke /tasks (yang akan mengarahkan ke /login jika belum autentikasi)
export default function Home() {
  redirect('/tasks');
}
