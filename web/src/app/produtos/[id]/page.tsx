"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { catalogApi, type Product } from "@/lib/api";
import { ChevronLeft } from "lucide-react";

function RelatedProductCard({ product }: { product: Product }) {
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
        <p className="text-2xl font-medium text-foreground">R$ {product.price}</p>
        <p className="mt-1 text-sm font-semibold text-ml-green">Frete grátis</p>
        <h3 className="mt-2 text-sm font-normal leading-tight text-muted line-clamp-2">
          {product.title}
        </h3>
      </div>
    </Link>
  );
}

export default function ProductDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params?.id as string | undefined;

  const [product, setProduct] = useState<Product | null>(null);
  const [related, setRelated] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    setError("");
    Promise.all([
      catalogApi.product(id),
      catalogApi.relatedProducts(id).then((res) => res.products),
    ])
      .then(([p, relatedProducts]) => {
        setProduct(p);
        setRelated(relatedProducts);
        setSelectedImageIndex(0);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Erro"))
      .finally(() => setLoading(false));
  }, [id]);

  if (!id) {
    return null;
  }

  if (loading) {
    return (
      <div className="pb-12 pt-8">
        <div className="mx-auto max-w-6xl px-4">
          <div className="h-8 w-48 animate-pulse rounded bg-gray-200" />
          <div className="mt-8 grid gap-8 md:grid-cols-2">
            <div className="aspect-square animate-pulse rounded-md bg-gray-200" />
            <div className="space-y-4">
              <div className="h-10 w-3/4 animate-pulse rounded bg-gray-200" />
              <div className="h-8 w-24 animate-pulse rounded bg-gray-200" />
              <div className="h-32 animate-pulse rounded bg-gray-200" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="pb-12 pt-8">
        <div className="mx-auto max-w-6xl px-4">
          <p className="text-red-600">{error || "Produto não encontrado."}</p>
          <Link href="/produtos" className="mt-4 inline-block text-ml-blue hover:underline">
            Voltar aos produtos
          </Link>
        </div>
      </div>
    );
  }

  const images = product.images?.length ? product.images : ["/placeholder-product.svg"];
  const currentImage = images[selectedImageIndex] ?? images[0];

  return (
    <div className="pb-12 pt-8">
      <div className="mx-auto max-w-6xl px-4">
        <button
          type="button"
          onClick={() => router.back()}
          className="mb-6 flex items-center gap-1 text-sm text-muted hover:text-foreground"
        >
          <ChevronLeft className="h-4 w-4" />
          Voltar
        </button>

        <div className="grid gap-8 md:grid-cols-2">
          {/* Gallery */}
          <div className="space-y-3">
            <div className="flex aspect-square items-center justify-center overflow-hidden rounded-lg border border-gray-200 bg-white p-4">
              {currentImage ? (
                <img
                  src={currentImage}
                  alt={product.title}
                  className="h-full w-full object-contain"
                />
              ) : (
                <span className="text-6xl text-gray-300">📦</span>
              )}
            </div>
            {images.length > 1 && (
              <div className="flex gap-2 overflow-x-auto pb-2">
                {images.map((src, index) => (
                  <button
                    key={index}
                    type="button"
                    onClick={() => setSelectedImageIndex(index)}
                    className={`h-16 w-16 shrink-0 overflow-hidden rounded border-2 bg-white ${
                      selectedImageIndex === index
                        ? "border-ml-blue"
                        : "border-gray-200 hover:border-gray-300"
                    }`}
                  >
                    {src ? (
                      <img src={src} alt="" className="h-full w-full object-contain" />
                    ) : (
                      <span className="text-2xl text-gray-300">📦</span>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Info */}
          <div>
            <h1 className="text-2xl font-semibold text-foreground md:text-3xl">
              {product.title}
            </h1>
            <p className="mt-4 text-4xl font-medium text-foreground">R$ {product.price}</p>
            <p className="mt-1 text-sm font-semibold text-ml-green">Frete grátis</p>
            <div className="mt-6">
              <h2 className="text-sm font-semibold text-muted">Descrição</h2>
              <p className="mt-2 whitespace-pre-wrap text-foreground">{product.description}</p>
            </div>
          </div>
        </div>

        {/* Related products */}
        {related.length > 0 && (
          <section className="mt-16">
            <h2 className="mb-6 text-xl font-semibold text-foreground">
              Produtos relacionados
            </h2>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
              {related.map((relatedProduct) => (
                <RelatedProductCard key={relatedProduct.id} product={relatedProduct} />
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  );
}
