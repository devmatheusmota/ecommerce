"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { catalogApi, type Product } from "@/lib/api";

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

export default function HomePage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    catalogApi
      .products({ limit: 10 })
      .then((res) => setProducts(res.products))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="pb-12 pt-8">
      {/* Hero Banner Placeholder */}
      <div className="mb-12 overflow-hidden rounded-md bg-gradient-to-r from-ml-blue to-blue-400">
        <div className="px-8 py-16 text-white md:py-24">
          <h1 className="text-3xl font-bold md:text-5xl">
            As melhores ofertas<br />estão aqui
          </h1>
          <p className="mt-4 text-lg text-blue-100 md:text-xl">
            Compre com segurança e receba no dia seguinte.
          </p>
          <Link
            href="/produtos"
            className="mt-8 inline-block rounded-md bg-ml-yellow px-8 py-3 font-semibold text-foreground transition-colors hover:bg-yellow-400"
          >
            Ver todas as ofertas
          </Link>
        </div>
      </div>

      {/* Featured Products */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-medium text-muted">Baseado na sua última visita</h2>
        <Link href="/produtos" className="text-sm font-medium text-ml-blue hover:text-ml-blue-hover">
          Ver histórico
        </Link>
      </div>

      {loading ? (
        <div className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-80 animate-pulse rounded-md bg-white shadow-sm"></div>
          ))}
        </div>
      ) : (
        <div className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}
