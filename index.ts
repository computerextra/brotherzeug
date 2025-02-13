import { readdir, appendFile, rm, readFile } from "fs/promises";
import { join } from "path";
import { PdfReader } from "pdfreader";
import path from "path";
import { fileURLToPath } from "url";
import { existsSync } from "fs";

const __filename = fileURLToPath(import.meta.url); // get the resolved path to the file
const __dirname = path.dirname(__filename); // get the name of the directory

const filePath = join(__dirname, "/files");

const eanRegex = /[0-9]{13}/;
const priceRegex = /^[0-9]*,[0-9]{2}/;

type Product = {
  ean: string;
  uvp?: string;
  sage?: string;
  preis?: string;
};

const all: Product[] = [];

async function readFiles() {
  const files = await readdir(filePath);

  let products: Product[] = [];

  files.forEach(async (file) => {
    if (file != ".gitkeep") {
      new PdfReader().parseFileItems(
        join(filePath, file),
        async (err, item) => {
          if (err) console.error(err);
          else if (!item) {
            products.forEach(async (x) => {
              await appendFile(join(__dirname, "temp.txt"), JSON.stringify(x));
            });
            all.push(...products);

            products = [];
          } else if (item.text) {
            // Hier ist der richtige Text

            if (item.text.match(priceRegex)) {
              products.push({
                preis: item.text,
                ean: "none",
              });
            }
            if (item.text.match(eanRegex)) {
              const i = products.find((x) => x.ean == "none");
              if (i) {
                i.ean = item.text;
              }
            }
          }
        }
      );
    }
  });
}

async function main() {
  if (existsSync(join(__dirname, "temp.txt"))) {
    await rm(join(__dirname, "temp.txt"));
  }
  const promise1 = Promise.resolve(readFiles());

  // TODO: Read temp.txt file and search SAGE if product exist
  Promise.all([promise1]).then(() => {
    console.log(all);
  });
}

void main();
