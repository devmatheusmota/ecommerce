"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { catalogApi, type Product } from "@/lib/api";
import { ShieldCheck, Truck, Undo2 } from "lucide-react";

export default function ProductDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [activeImage, setActiveImage] = useState(0);

  useEffect(() => {
    catalogApi
      .product(id)
      .then(setProduct)
      .catch((err) => setError(err instanceof Error ? err.message : "Erro"))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return (
      <div className="pt-12">
        <div className="h-96 animate-pulse rounded-md bg-white shadow-sm"></div>
      </div>
    );
  }
  if (error || !product) {
    return (
      <div className="pt-12 text-center">
        <p className="text-lg text-red-600">{error || "Produto não encontrado"}</p>
        <Link href="/produtos" className="mt-4 inline-block font-medium text-ml-blue hover:underline">
          Voltar aos produtos
        </Link>
      </div>
    );
  }

  const images = product.images?.length ? product.images : ["/placeholder-product.svg"];

  return (
    <div className="pb-16 pt-8">
      <div className="rounded-md bg-white p-6 shadow-sm md:p-8">
        <div className="grid gap-8 lg:grid-cols-12">
          
          {/* Left Column: Images */}
          <div className="lg:col-span-7 flex gap-4">
            {/* Thumbnails */}
            <div className="flex w-16 flex-col gap-2">
              {images.map((img, idx) => (
                <button
                  key={idx}
                  onMouseEnter={() => setActiveImage(idx)}
                  className={`aspect-square overflow-hidden rounded-md border-2 p-1 ${activeImage === idx ? "border-ml-blue" : "border-transparent hover:border-gray-300"}`}
                >
                  <img src={img} alt="" className="h-full w-full object-contain" />
                </button>
              ))}
            </div>
            {/* Main Image */}
            <div className="flex flex-1 items-center justify-center overflow-hidden rounded-md p-4">
              <img
                src={images[activeImage]}
                alt={product.title}
                className="max-h-[500px] w-full object-contain"
              />
            </div>
          </div>

          {/* Right Column: Info & Buy Box */}
          <div className="lg:col-span-5 flex flex-col rounded-md border border-gray-200 p-6">
            <p className="text-sm text-muted">Novo | 123 vendidos</p>
            <h1 className="mt-2 text-2xl font-bold text-foreground sm:text-3xl">
              {product.title}
            </h1>
            
            <div className="mt-6">
              <p className="text-4xl font-light text-foreground">
                R$ {product.price}
              </p>
              <p className="mt-1 text-sm text-foreground">
                em <span className="font-medium text-ml-green">10x R$ {(Number(product.price) / 10).toFixed(2)} sem juros</span>
              </p>
            </div>

            <div className="mt-6 space-y-4">
              <div className="flex items-start gap-3 text-sm">
                <Truck className="mt-0.5 h-5 w-5 text-ml-green shrink-0" />
                <div>
                  <p className="font-medium text-ml-green">Chegará grátis amanhã</p>
                  <p className="text-muted">Comprando dentro das próximas 2 h</p>
                </div>
              </div>
              <div className="flex items-start gap-3 text-sm">
                <Undo2 className="mt-0.5 h-5 w-5 text-ml-blue shrink-0" />
                <div>
                  <p className="font-medium text-ml-blue">Devolução grátis</p>
                  <p className="text-muted">Você tem 30 dias a partir do recebimento</p>
                </div>
              </div>
              <div className="flex items-start gap-3 text-sm">
                <ShieldCheck className="mt-0.5 h-5 w-5 text-muted shrink-0" />
                <div>
                  <p className="font-medium text-ml-blue">Compra Garantida</p>
                  <p className="text-muted">Receba o produto que está esperando ou devolvemos o dinheiro</p>
                </div>
              </div>
            </div>

            <div className="mt-8 flex flex-col gap-3">
              <button className="rounded-md bg-ml-blue py-4 font-semibold text-white transition-colors hover:bg-ml-blue-hover">
                Comprar agora
              </button>
              <button className="rounded-md bg-blue-100 py-4 font-semibold text-ml-blue transition-colors hover:bg-blue-200">
                Adicionar ao carrinho
              </button>
            </div>
          </div>
        </div>

        {/* Description Section */}
        <div className="mt-16 border-t border-gray-100 pt-12">
          <h2 className="text-2xl font-normal text-foreground">Descrição do produto</h2>
          <p className="mt-6 whitespace-pre-wrap text-lg text-muted leading-relaxed">
            {product.description}
          </p>
        </div>
      </div>
    </div>
  );
}
