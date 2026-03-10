"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { catalogApi, usersApi } from "@/lib/api";
import { getToken } from "@/lib/auth";

type Category = { id: string; name: string; slug: string };

export default function NovoProdutoPage() {
  const router = useRouter();
  const [categories, setCategories] = useState<Category[]>([]);
  const [sellerId, setSellerId] = useState("");
  const [form, setForm] = useState({
    category_id: "",
    title: "",
    description: "",
    price: "",
    images: "",
  });
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    catalogApi
      .categories()
      .then(setCategories)
      .catch(() => setCategories([]))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    const token = getToken();
    if (token) {
      usersApi
        .me(token)
        .then((user) => setSellerId(user.id))
        .catch(() => {});
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!sellerId.trim()) {
      setError("Faça login para publicar produtos (ou informe o ID do vendedor)");
      return;
    }
    if (!form.category_id.trim()) {
      setError("Selecione uma categoria");
      return;
    }
    setSubmitting(true);
    try {
      const product = await catalogApi.createProduct({
        seller_id: sellerId,
        category_id: form.category_id,
        title: form.title.trim(),
        description: form.description.trim(),
        price: form.price.trim().replace(",", "."),
        images: form.images
          ? form.images
              .split("\n")
              .map((url) => url.trim())
              .filter(Boolean)
          : undefined,
      });
      router.push(`/produtos/${product.id}`);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao criar produto");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="pt-12">
        <p className="text-[var(--muted)]">Carregando categorias...</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-xl pt-12">
      <Link
        href="/produtos"
        className="mb-6 inline-block text-sm text-[var(--muted)] hover:text-[var(--accent)]"
      >
        ← Voltar aos produtos
      </Link>
      <h1 className="font-serif text-2xl font-bold">Anunciar produto</h1>
      <p className="mt-2 text-[var(--muted)]">
        Preencha os dados para criar seu anúncio
      </p>
      <form onSubmit={handleSubmit} className="mt-8 space-y-4">
        <div>
          <label htmlFor="category_id" className="block text-sm font-medium">
            Categoria *
          </label>
          <select
            id="category_id"
            required
            value={form.category_id}
            onChange={(e) =>
              setForm((f) => ({ ...f, category_id: e.target.value }))
            }
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          >
            <option value="">Selecione...</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="title" className="block text-sm font-medium">
            Título *
          </label>
          <input
            id="title"
            type="text"
            required
            maxLength={500}
            value={form.title}
            onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
            placeholder="Ex: Smartphone Samsung Galaxy..."
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="description" className="block text-sm font-medium">
            Descrição *
          </label>
          <textarea
            id="description"
            required
            rows={4}
            value={form.description}
            onChange={(e) =>
              setForm((f) => ({ ...f, description: e.target.value }))
            }
            placeholder="Descreva o produto..."
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="price" className="block text-sm font-medium">
            Preço (R$) *
          </label>
          <input
            id="price"
            type="text"
            required
            value={form.price}
            onChange={(e) => setForm((f) => ({ ...f, price: e.target.value }))}
            placeholder="99.90"
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        <div>
          <label htmlFor="images" className="block text-sm font-medium">
            URLs das imagens (uma por linha)
          </label>
          <textarea
            id="images"
            rows={3}
            value={form.images}
            onChange={(e) => setForm((f) => ({ ...f, images: e.target.value }))}
            placeholder="https://exemplo.com/img1.jpg"
            className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
          />
        </div>
        {!sellerId && (
          <div>
            <label htmlFor="seller_id" className="block text-sm font-medium">
              ID do vendedor (UUID)
            </label>
            <input
              id="seller_id"
              type="text"
              value={sellerId}
              onChange={(e) => setSellerId(e.target.value)}
              placeholder="00000000-0000-0000-0000-000000000001"
              className="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2.5 outline-none focus:border-ml-blue focus:ring-1 focus:ring-ml-blue"
            />
            <p className="mt-1 text-xs text-[var(--muted)]">
              Faça login para usar seu ID automaticamente
            </p>
          </div>
        )}
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          disabled={submitting}
          className="w-full rounded-md bg-ml-blue py-3 font-semibold text-white transition-colors hover:bg-ml-blue-hover disabled:opacity-60"
        >
          {submitting ? "Criando..." : "Criar produto"}
        </button>
      </form>
    </div>
  );
}
