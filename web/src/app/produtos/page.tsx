"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { catalogApi, type Product, type CategoryTreeNode } from "@/lib/api";

function ProductCard({ product }: { product: Product }) {
  const imageUrl = product.images?.[0] || "/placeholder-product.svg";

  return (
    <Link
      href={`/produtos/${product.id}`}
      className="group flex flex-col overflow-hidden rounded-md bg-white shadow-sm transition-shadow hover:shadow-md"
    >
      <div className="flex aspect-square items-center justify-center overflow-hidden border-b border-gray-100 bg-white p-4">
        {imageUrl.startsWith("http") ? (
          <img
            src={imageUrl}
            alt={product.title}
            className="h-full w-full object-contain transition-transform group-hover:scale-105"
          />
        ) : (
          <span className="text-4xl text-gray-300">📦</span>
        )}
      </div>
      <div className="flex flex-1 flex-col p-4">
        <p className="text-2xl font-medium text-foreground">
          R$ {product.price}
        </p>
        <p className="mt-1 text-sm font-semibold text-ml-green">
          Frete grátis
        </p>
        <h3 className="mt-2 text-sm font-normal leading-tight text-muted line-clamp-2">
          {product.title}
        </h3>
      </div>
    </Link>
  );
}

function ProdutosContent() {
  const searchParams = useSearchParams();
  const categoryId = searchParams.get("category_id") || undefined;
  
  const [result, setResult] = useState<{ products: Product[]; total: number } | null>(null);
  const [categories, setCategories] = useState<CategoryTreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    catalogApi.categoriesTree().then(setCategories).catch(() => {});
  }, []);

  useEffect(() => {
    setLoading(true);
    catalogApi
      .products({ category_id: categoryId, limit: 24, offset: 0 })
      .then(setResult)
      .catch((err) => setError(err instanceof Error ? err.message : "Erro"))
      .finally(() => setLoading(false));
  }, [categoryId]);

  return (
    <div className="flex flex-col gap-8 pb-12 pt-8 md:flex-row">
      {/* Sidebar - Categories */}
      <aside className="w-full shrink-0 md:w-64">
        <div className="rounded-md bg-white p-6 shadow-sm">
          <h2 className="font-semibold text-foreground">Categorias</h2>
          <ul className="mt-4 space-y-2 text-sm text-muted">
            <li>
              <Link
                href="/produtos"
                className={`block py-1 hover:text-ml-blue ${!categoryId ? "font-medium text-ml-blue" : ""}`}
              >
                Todas as categorias
              </Link>
            </li>
            {categories.map((cat) => (
              <li key={cat.id}>
                <Link
                  href={`/produtos?category_id=${cat.id}`}
                  className={`block py-1 hover:text-ml-blue ${categoryId === cat.id ? "font-medium text-ml-blue" : ""}`}
                >
                  {cat.name}
                </Link>
                {cat.children && cat.children.length > 0 && (
                  <ul className="ml-4 mt-1 space-y-1 border-l border-gray-200 pl-4">
                    {cat.children.map((child) => (
                      <li key={child.id}>
                        <Link
                          href={`/produtos?category_id=${child.id}`}
                          className={`block py-1 hover:text-ml-blue ${categoryId === child.id ? "font-medium text-ml-blue" : ""}`}
                        >
                          {child.name}
                        </Link>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            ))}
          </ul>
        </div>
      </aside>

      {/* Main Content - Products */}
      <main className="flex-1">
        <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold text-foreground">
              {categoryId ? "Resultados da categoria" : "Todos os produtos"}
            </h1>
            <p className="text-sm text-muted">
              {result?.total || 0} {(result?.total || 0) === 1 ? "resultado" : "resultados"}
            </p>
          </div>
          <Link
            href="/produtos/novo"
            className="rounded-md bg-ml-blue px-5 py-2.5 font-medium text-white transition-colors hover:bg-ml-blue-hover"
          >
            Anunciar produto
          </Link>
        </div>

        {loading ? (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 xl:grid-cols-4">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="h-80 animate-pulse rounded-md bg-white shadow-sm"></div>
            ))}
          </div>
        ) : error ? (
          <div className="rounded-md bg-red-50 p-4 text-red-600">{error}</div>
        ) : result?.products.length === 0 ? (
          <div className="rounded-md bg-white p-12 text-center shadow-sm">
            <p className="text-lg text-muted">Nenhum produto encontrado.</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 xl:grid-cols-4">
            {result?.products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}

export default function ProdutosPage() {
  return (
    <Suspense fallback={<p className="pt-12 text-muted">Carregando...</p>}>
      <ProdutosContent />
    </Suspense>
  );
}
