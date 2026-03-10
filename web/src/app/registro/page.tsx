"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { usersApi } from "@/lib/api";
import { setToken } from "@/lib/auth";

export default function RegisterPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    email: "",
    password: "",
    name: "",
    phone: "",
    cpf: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await usersApi.register(form);
      const { token } = await usersApi.login({
        email: form.email,
        password: form.password,
      });
      setToken(token);
      router.push("/perfil");
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao cadastrar");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-md pt-12">
      <h1 className="font-serif text-2xl font-bold">Criar conta</h1>
      <p className="mt-2 text-[var(--muted)]">
        Preencha os dados para se cadastrar
      </p>
      <form onSubmit={handleSubmit} className="mt-8 space-y-4">
        <div>
          <label htmlFor="name" className="block text-sm font-medium">
            Nome
          </label>
          <input
            id="name"
            type="text"
            required
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="email" className="block text-sm font-medium">
            E-mail
          </label>
          <input
            id="email"
            type="email"
            required
            value={form.email}
            onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="password" className="block text-sm font-medium">
            Senha (mín. 6 caracteres)
          </label>
          <input
            id="password"
            type="password"
            required
            minLength={6}
            value={form.password}
            onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="phone" className="block text-sm font-medium">
            Telefone
          </label>
          <input
            id="phone"
            type="tel"
            required
            value={form.phone}
            onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))}
            placeholder="11999999999"
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="cpf" className="block text-sm font-medium">
            CPF
          </label>
          <input
            id="cpf"
            type="text"
            required
            value={form.cpf}
            onChange={(e) => setForm((f) => ({ ...f, cpf: e.target.value }))}
            placeholder="12345678900"
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        {error && (
          <p className="text-sm text-red-600">{error}</p>
        )}
        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-md bg-ml-blue py-3 font-semibold text-white transition-colors hover:bg-ml-blue-hover disabled:opacity-60"
        >
          {loading ? "Cadastrando..." : "Cadastrar"}
        </button>
      </form>
      <p className="mt-6 text-center text-sm text-[var(--muted)]">
        Já tem conta?{" "}
        <Link href="/login" className="font-medium text-ml-blue hover:underline">
          Entrar
        </Link>
      </p>
    </div>
  );
}
