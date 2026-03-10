"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { usersApi } from "@/lib/api";
import { getToken } from "@/lib/auth";

type User = {
  id: string;
  email: string;
  name: string;
  phone: string;
  cpf: string;
};

export default function ProfilePage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push("/login");
      return;
    }
    usersApi
      .me(token)
      .then(setUser)
      .catch((err) => {
        setError(err instanceof Error ? err.message : "Erro ao carregar perfil");
        if (err instanceof Error && err.message.includes("identity")) {
          router.push("/login");
        }
      })
      .finally(() => setLoading(false));
  }, [router]);

  if (loading) {
    return (
      <div className="pt-12 text-center">
        <p className="text-[var(--muted)]">Carregando perfil...</p>
      </div>
    );
  }

  if (error && !user) {
    return (
      <div className="pt-12">
        <p className="text-red-600">{error}</p>
        <Link href="/login" className="mt-4 inline-block text-ml-blue hover:underline">
          Ir para login
        </Link>
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="mx-auto max-w-md pt-12">
      <h1 className="font-serif text-2xl font-bold">Meu perfil</h1>
      <div className="mt-8 rounded-md border border-gray-200 bg-white p-6 shadow-sm">
        <dl className="space-y-3">
          <div>
            <dt className="text-sm text-[var(--muted)]">Nome</dt>
            <dd className="font-medium">{user.name}</dd>
          </div>
          <div>
            <dt className="text-sm text-[var(--muted)]">E-mail</dt>
            <dd className="font-medium">{user.email}</dd>
          </div>
          <div>
            <dt className="text-sm text-[var(--muted)]">Telefone</dt>
            <dd className="font-medium">{user.phone || "—"}</dd>
          </div>
          <div>
            <dt className="text-sm text-[var(--muted)]">CPF</dt>
            <dd className="font-medium">{user.cpf || "—"}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
}
