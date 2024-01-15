import { FormEvent, FunctionComponent, useState } from "react";

import { Button } from "@/app/components/ui/button";
import LoadingDots from "@/app/components/ui/loadingdots";
import { useRouter } from "next/navigation";

const PairCharger = () => {
  const [isLoading, setIsLoading] = useState(false);

  const [chargerID, setChargerID] = useState("");
  const router = useRouter();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    router.push(`/charge/${formData.get("chargerID")}`);
  };

  return (
    <div className="flex flex-col justify-between h-stretch">
      <div className="flex flex-col self-stretch gap-16 p-6 mx-auto md align-">
        {/* start banner */}
        <div className="flex flex-col items-center justify-center gap-3 px-4 py-0 ">
          <div className="">
            <img
              className="object-scale-down pb-12 h-52 md:h-56"
              alt="Exercise"
              src="/homework.png"
            />
          </div>
          <div className="flex flex-col gap-6">
            <p className="text-2xl font-normal text-center text-slate-light-12">
              Welcome to the
              <br />
              DejaBlue software interview
            </p>
            <p className="text-base text-center text-slate-light-11">
              Please start by pairing your charger.
            </p>
          </div>
        </div>

        {/* input section */}
        <div className="flex flex-col items-center gap-16">
          <form onSubmit={handleSubmit} className="flex flex-col w-full gap-8">
            <div className="rounded-md px-3 pb-1.5 pt-2.5 shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-blue-500">
              <label
                htmlFor="chargerID"
                className="block pb-1 text-xs font-medium text-slate-light-11"
              >
                Charger code
              </label>
              <input
                type="text"
                name="chargerID"
                id="chargerID"
                required={true}
                className="block w-full p-0 border-0 text-slate-light-12 placeholder:text-slate-light-11 focus:ring-0 sm:text-sm sm:leading-6"
                placeholder="DEF"
                onInput={(e) => {
                  const inputElement = e.target as HTMLInputElement;
                  inputElement.setCustomValidity("");
                }}
                onInvalid={(e) => {
                  const inputElement = e.target as HTMLInputElement;
                  inputElement.setCustomValidity("Please enter a charger code");
                }}
              />
            </div>

            <Button
              type="submit"
              className="w-48 mx-auto text-white bg-black rounded-full h-14 hover:bg-blue-700"
            >
              {isLoading ? <LoadingDots style={"small"} /> : "Pair Charger"}
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default PairCharger;
