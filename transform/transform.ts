interface ComplexElement {
  element_id: number;
  another_data_field: string;
}

export interface TreeNode {
  id: number;
  some_data_field: string;
  children: TreeNode[];
  elements: ComplexElement[];
}

export interface WidgetTree {
  id: number;
  some_data_field: string;
  open: boolean;
  data: (ComplexElement | WidgetTree)[];
}

export function transform(tree: TreeNode): WidgetTree {
  let items: (ComplexElement | WidgetTree)[] = [];
  for (let child of tree.children) {
    items.push(transform(child));
  }
  items.concat(tree.elements);
  return {
    id: tree.id,
    some_data_field: tree.some_data_field,
    open: false,
    data: items,
  };
}
