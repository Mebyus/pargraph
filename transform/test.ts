import { transform, TreeNode, WidgetTree } from "transform";

interface TestCase {
  name: string;
  args: {
    tree: TreeNode;
  };
  want: WidgetTree;
}

function test_transform(): void {
  let testCases: TestCase[] = [
    {
      name: "no elements, no children",
      args: {
        tree: {
          id: 2,
          some_data_field: "some_data",
          children: [],
          elements: [],
        },
      },
      want: {
        id: 2,
        some_data_field: "some_data",
        open: true,
        data: [],
      },
    },
  ];
  for (let testCase of testCases) {
    console.log(testCase.name);
    console.log(transform(testCase.args.tree));
  }
}

test_transform();
