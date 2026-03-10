const API_BASE = ""; // Same origin via Next.js rewrites

type ApiResponse<T> = { data: T; meta?: { timestamp: string; version: string } };
type ApiError = { error: string };

async function request<T>(
  path: string,
  options: RequestInit & { token?: string } = {}
): Promise<T> {
  const { token, ...init } = options;
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...init.headers,
  };
  if (token) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${path}`, { ...init, headers });
  const json = await response.json();

  if (!response.ok) {
    const err = json as ApiError;
    throw new Error(err.error || response.statusText);
  }

  const wrapped = json as ApiResponse<T>;
  return wrapped.data ?? json;
}

// Users
export const usersApi = {
  register: (body: {
    email: string;
    password: string;
    name: string;
    phone: string;
    cpf: string;
  }) => request<{ id: string; email: string; name: string }>("/api/users/register", {
    method: "POST",
    body: JSON.stringify(body),
  }),

  login: (body: { email: string; password: string }) =>
    request<{ token: string; expire_at: string }>("/api/users/login", {
      method: "POST",
      body: JSON.stringify(body),
    }),

  me: (token: string) =>
    request<{ id: string; email: string; name: string; phone: string; cpf: string }>(
      "/api/users/me",
      { headers: { Authorization: `Bearer ${token}` } }
    ),

  updateMe: (token: string, body: { name?: string; phone?: string; cpf?: string }) =>
    request<{ id: string; email: string; name: string; phone: string; cpf: string }>(
      "/api/users/me",
      {
        method: "PATCH",
        body: JSON.stringify(body),
        headers: { Authorization: `Bearer ${token}` },
      }
    ),
};

// Catalog
export const catalogApi = {
  categoriesTree: () =>
    request<CategoryTreeNode[]>("/api/catalog/v1/categories/tree"),

  categories: (parentId?: string) => {
    const q = parentId ? `?parent_id=${parentId}` : "";
    return request<
      Array<{
        id: string;
        name: string;
        slug: string;
        parent_id?: string;
      }>
    >(`/api/catalog/v1/categories${q}`);
  },

  products: (params?: { seller_id?: string; category_id?: string; limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.seller_id) searchParams.set("seller_id", params.seller_id);
    if (params?.category_id) searchParams.set("category_id", params.category_id);
    if (params?.limit) searchParams.set("limit", String(params.limit));
    if (params?.offset) searchParams.set("offset", String(params.offset));
    const q = searchParams.toString() ? `?${searchParams.toString()}` : "";
    return request<{ products: Array<Product>; total: number }>(`/api/catalog/v1/products${q}`);
  },

  product: (id: string) =>
    request<Product>(`/api/catalog/v1/products/${id}`),

  relatedProducts: (productId: string) =>
    request<{ products: Array<Product> }>(`/api/catalog/v1/products/${productId}/related`),

  createProduct: (body: {
    seller_id: string;
    category_id: string;
    title: string;
    description: string;
    price: string;
    images?: string[];
  }) =>
    request<Product>("/api/catalog/v1/products", {
      method: "POST",
      body: JSON.stringify(body),
    }),
};

export type CategoryTreeNode = {
  id: string;
  name: string;
  slug: string;
  parent_id?: string;
  children: CategoryTreeNode[];
};

export type Product = {
  id: string;
  seller_id: string;
  category_id: string;
  title: string;
  description: string;
  price: string;
  images: string[];
};
