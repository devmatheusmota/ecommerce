"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { catalogApi, type Product, type CategoryTreeNode } from "@/lib/api";
import { ChevronLeft, ChevronRight } from "lucide-react";

// All divisible by 3 and 4 so the grid (3–4 cols) always has full rows
const PAGE_SIZE_OPTIONS = [12, 24, 36, 48, 96] as const;
const DEFAULT_PAGE_SIZE = 24;

function ProductCard({ product }: { product: Product }) {
  const imageUrl = product.images?.[0] || "/placeholder-product.svg";

  return (
    <Link
      href={`/produtos/${product.id}`}
      className="group flex flex-col overflow-hidden rounded-md bg-white shadow-sm transition-shadow hover:shadow-md"
    >
      <div className="flex aspect-square items-center justify-center overflow-hidden border-b border-gray-100 bg-white p-4">
        {imageUrl ? (
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

function parsePageSize(value: string | null): number {
  const n = value ? parseInt(value, 10) : NaN;
  return PAGE_SIZE_OPTIONS.includes(n as (typeof PAGE_SIZE_OPTIONS)[number])
    ? (n as (typeof PAGE_SIZE_OPTIONS)[number])
    : DEFAULT_PAGE_SIZE;
}

function ProdutosContent() {
  const searchParams = useSearchParams();
  const categoryId = searchParams.get("category_id") || undefined;
  const page = Math.max(1, parseInt(searchParams.get("page") || "1", 10) || 1);
  const pageSize = parsePageSize(searchParams.get("per_page"));

  const [result, setResult] = useState<{ products: Product[]; total: number } | null>(null);
  const [categories, setCategories] = useState<CategoryTreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    catalogApi.categoriesTree().then(setCategories).catch(() => {});
  }, []);

  useEffect(() => {
    setLoading(true);
    const offset = (page - 1) * pageSize;
    catalogApi
      .products({ category_id: categoryId, limit: pageSize, offset })
      .then(setResult)
      .catch((err) => setError(err instanceof Error ? err.message : "Erro"))
      .finally(() => setLoading(false));
  }, [categoryId, page, pageSize]);

  const total = result?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const hasPrev = page > 1;
  const hasNext = page < totalPages;

  function buildPageUrl(nextPage: number) {
    const params = new URLSearchParams(searchParams.toString());
    if (nextPage <= 1) params.delete("page");
    else params.set("page", String(nextPage));
    const q = params.toString();
    return `/produtos${q ? `?${q}` : ""}`;
  }

  function buildPerPageUrl(newPerPage: number) {
    const params = new URLSearchParams(searchParams.toString());
    params.set("per_page", String(newPerPage));
    params.delete("page");
    const q = params.toString();
    return `/produtos${q ? `?${q}` : ""}`;
  }

  return (
    <div className="flex flex-col gap-8 pb-12 pt-8 md:flex-row">
      {/* Sidebar - Categories */}
      <aside className="w-full shrink-0 md:w-64">
        <div className="rounded-md bg-white p-6 shadow-sm">
          <h2 className="font-semibold text-foreground">Categorias</h2>
          <ul className="mt-4 space-y-2 text-sm text-muted">
            <li>
              <Link
                href={pageSize !== DEFAULT_PAGE_SIZE ? `/produtos?per_page=${pageSize}` : "/produtos"}
                className={`block py-1 hover:text-ml-blue ${!categoryId ? "font-medium text-ml-blue" : ""}`}
              >
                Todas as categorias
              </Link>
            </li>
            {categories.map((cat) => (
              <li key={cat.id}>
                <Link
                  href={
                    pageSize !== DEFAULT_PAGE_SIZE
                      ? `/produtos?category_id=${cat.id}&per_page=${pageSize}`
                      : `/produtos?category_id=${cat.id}`
                  }
                  className={`block py-1 hover:text-ml-blue ${categoryId === cat.id ? "font-medium text-ml-blue" : ""}`}
                >
                  {cat.name}
                </Link>
                {cat.children && cat.children.length > 0 && (
                  <ul className="ml-4 mt-1 space-y-1 border-l border-gray-200 pl-4">
                    {cat.children.map((child) => (
                      <li key={child.id}>
                        <Link
                          href={
                            pageSize !== DEFAULT_PAGE_SIZE
                              ? `/produtos?category_id=${child.id}&per_page=${pageSize}`
                              : `/produtos?category_id=${child.id}`
                          }
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
              {total} {total === 1 ? "resultado" : "resultados"}
              {totalPages > 1 && ` · Página ${page} de ${totalPages}`}
            </p>
          </div>
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center gap-2 text-sm text-muted">
              <span>Exibir por página:</span>
              <span className="flex gap-1">
                {PAGE_SIZE_OPTIONS.map((size) =>
                  size === pageSize ? (
                    <span
                      key={size}
                      className="flex h-8 min-w-[2rem] items-center justify-center rounded border border-ml-blue bg-ml-blue px-2 font-medium text-white"
                    >
                      {size}
                    </span>
                  ) : (
                    <Link
                      key={size}
                      href={buildPerPageUrl(size)}
                      className="flex h-8 min-w-[2rem] items-center justify-center rounded border border-gray-300 bg-white px-2 font-medium text-foreground transition-colors hover:bg-gray-50"
                    >
                      {size}
                    </Link>
                  )
                )}
              </span>
            </div>
            <Link
              href="/produtos/novo"
              className="rounded-md bg-ml-blue px-5 py-2.5 font-medium text-white transition-colors hover:bg-ml-blue-hover"
            >
              Anunciar produto
            </Link>
          </div>
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
          <>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 xl:grid-cols-4">
              {result?.products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <nav className="mt-10 flex flex-wrap items-center justify-center gap-2" aria-label="Paginação">
                {hasPrev ? (
                  <Link
                    href={buildPageUrl(page - 1)}
                    className="flex items-center gap-1 rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-gray-50"
                  >
                    <ChevronLeft className="h-4 w-4" />
                    Anterior
                  </Link>
                ) : (
                  <span className="flex cursor-not-allowed items-center gap-1 rounded-md border border-gray-200 bg-gray-50 px-4 py-2 text-sm text-muted">
                    <ChevronLeft className="h-4 w-4" />
                    Anterior
                  </span>
                )}

                <div className="flex items-center gap-1">
                  {Array.from({ length: totalPages }, (_, i) => i + 1)
                    .filter((p) => {
                      if (totalPages <= 7) return true;
                      if (p === 1 || p === totalPages) return true;
                      if (Math.abs(p - page) <= 1) return true;
                      return false;
                    })
                    .map((p, idx, arr) => {
                      const showEllipsisBefore = idx > 0 && arr[idx - 1] !== p - 1;
                      return (
                        <span key={p} className="flex items-center gap-1">
                          {showEllipsisBefore && <span className="px-2 text-muted">…</span>}
                          {p === page ? (
                            <span className="flex h-9 w-9 items-center justify-center rounded-md bg-ml-blue font-medium text-white">
                              {p}
                            </span>
                          ) : (
                            <Link
                              href={buildPageUrl(p)}
                              className="flex h-9 w-9 items-center justify-center rounded-md border border-gray-300 bg-white text-sm font-medium text-foreground transition-colors hover:bg-gray-50"
                            >
                              {p}
                            </Link>
                          )}
                        </span>
                      );
                    })}
                </div>

                {hasNext ? (
                  <Link
                    href={buildPageUrl(page + 1)}
                    className="flex items-center gap-1 rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-gray-50"
                  >
                    Próxima
                    <ChevronRight className="h-4 w-4" />
                  </Link>
                ) : (
                  <span className="flex cursor-not-allowed items-center gap-1 rounded-md border border-gray-200 bg-gray-50 px-4 py-2 text-sm text-muted">
                    Próxima
                    <ChevronRight className="h-4 w-4" />
                  </span>
                )}
              </nav>
            )}
          </>
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
