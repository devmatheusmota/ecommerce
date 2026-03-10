"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { catalogApi, type CategoryTreeNode } from "@/lib/api";
import { ChevronRight } from "lucide-react";

function CategoryTree({ categories }: { categories: CategoryTreeNode[] }) {
  if (categories.length === 0) {
    return (
      <div className="rounded-md bg-white p-12 text-center shadow-sm">
        <p className="text-lg text-muted">Nenhuma categoria cadastrada ainda.</p>
      </div>
    );
  }
  return (
    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {categories.map((cat) => (
        <div key={cat.id} className="rounded-md bg-white p-6 shadow-sm transition-shadow hover:shadow-md">
          <Link
            href={`/produtos?category_id=${cat.id}`}
            className="flex items-center justify-between font-semibold text-foreground hover:text-ml-blue"
          >
            <span className="text-lg">{cat.name}</span>
            <ChevronRight className="h-5 w-5 text-gray-400" />
          </Link>
          
          {cat.children && cat.children.length > 0 && (
            <ul className="mt-4 space-y-3">
              {cat.children.map((child) => (
                <li key={child.id}>
                  <Link
                    href={`/produtos?category_id=${child.id}`}
                    className="text-sm text-muted hover:text-ml-blue"
                  >
                    {child.name}
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </div>
      ))}
    </div>
  );
}

export default function CategoriasPage() {
  const [categories, setCategories] = useState<CategoryTreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    catalogApi
      .categoriesTree()
      .then(setCategories)
      .catch((err) => setError(err instanceof Error ? err.message : "Erro"))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="pb-12 pt-8">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold text-foreground">Categorias para comprar e vender</h1>
      </div>
      
      {loading && (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="h-48 animate-pulse rounded-md bg-white shadow-sm"></div>
          ))}
        </div>
      )}
      
      {error && (
        <div className="rounded-md bg-red-50 p-4 text-red-600">{error}</div>
      )}
      
      {!loading && !error && (
        <CategoryTree categories={categories} />
      )}
    </div>
  );
}
